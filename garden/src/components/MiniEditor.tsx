import { createSignal, onMount, onCleanup, Show } from "solid-js";
import { createEditor, $getSelection, $isRangeSelection, $isElementNode, $createTextNode, $createParagraphNode, $insertNodes, $getRoot, COMMAND_PRIORITY_LOW } from "lexical";
import type { LexicalEditor, TextFormatType } from "lexical";
import { registerRichText } from "@lexical/rich-text";
import { $generateHtmlFromNodes, $generateNodesFromDOM } from "@lexical/html";
import { LinkNode, $createLinkNode, $isLinkNode, TOGGLE_LINK_COMMAND, $toggleLink } from "@lexical/link";
import { registerHistory, createEmptyHistoryState } from "@lexical/history";
import {
  IconBold, IconBoldOff,
  IconItalic,
  IconUnderline,
  IconEyeOff,
  IconLink, IconLinkOff,
  IconMoodSmile,
} from "@tabler/icons-solidjs";
import "emoji-picker-element";

interface MiniEditorProps {
  onHtml: (html: string) => void;
  initialHtml?: string;
}

export default function MiniEditor(props: MiniEditorProps) {
  let editorRef!: HTMLDivElement;
  let containerRef!: HTMLDivElement;
  let linkUrlRef!: HTMLInputElement;
  let editorInstance: LexicalEditor;

  const [bold, setBold] = createSignal(false);
  const [italic, setItalic] = createSignal(false);
  const [underline, setUnderline] = createSignal(false);
  const [spoiler, setSpoiler] = createSignal(false);
  const [link, setLink] = createSignal(false);
  const [showLinkInput, setShowLinkInput] = createSignal(false);
  const [linkUrl, setLinkUrl] = createSignal("");
  const [linkText, setLinkText] = createSignal("");
  const [hasSelection, setHasSelection] = createSignal(false);
  const [showEmoji, setShowEmoji] = createSignal(false);
  const [emojiPos, setEmojiPos] = createSignal({ top: 0, right: 0 });

  function ensureSelection() {
    let selection = $getSelection();
    if (!$isRangeSelection(selection)) {
      const root = $getRoot();
      const firstChild = root.getFirstChild();
      if ($isElementNode(firstChild)) {
        selection = firstChild.select(0, 0);
      } else {
        const paragraph = $createParagraphNode();
        root.append(paragraph);
        selection = paragraph.select(0, 0);
      }
    }
    return selection;
  }

  function ensureFocus() {
    editorInstance.update(() => { ensureSelection(); });
    editorInstance.focus();
  }

  function formatText(format: TextFormatType) {
    editorInstance.update(() => {
      const selection = ensureSelection();
      if ($isRangeSelection(selection)) {
        selection.formatText(format);
      }
    });
    editorInstance.focus();
  }

  function openLinkInput() {
    setShowEmoji(false);
    if (link()) {
      ensureFocus();
      editorInstance.dispatchCommand(TOGGLE_LINK_COMMAND, null);
      return;
    }

    editorInstance.getEditorState().read(() => {
      const selection = $getSelection();
      if ($isRangeSelection(selection)) {
        const text = selection.getTextContent();
        setLinkText(text);
        setHasSelection(text.length > 0);
      } else {
        setLinkText("");
        setHasSelection(false);
      }
    });

    setLinkUrl("");
    setShowLinkInput(true);
    requestAnimationFrame(() => linkUrlRef?.focus());
  }

  function submitLink() {
    const url = linkUrl().trim();
    if (!url) return;

    const text = linkText().trim();

    ensureFocus();
    editorInstance.update(() => {
      const selection = $getSelection();
      if ($isRangeSelection(selection) && hasSelection()) {
        $toggleLink(url);
      } else {
        const linkNode = $createLinkNode(url);
        linkNode.append($createTextNode(text || url));
        $insertNodes([linkNode]);
      }
    });

    setShowLinkInput(false);
    setLinkUrl("");
    setLinkText("");
  }

  function cancelLink() {
    setShowLinkInput(false);
    setLinkUrl("");
    setLinkText("");
    ensureFocus();
  }

  let emojiBtnRef!: HTMLButtonElement;

  function toggleEmojiPicker() {
    if (!showEmoji()) {
      const rect = emojiBtnRef.getBoundingClientRect();
      setEmojiPos({ top: rect.bottom + 2, right: window.innerWidth - rect.right });
    }
    setShowEmoji(!showEmoji());
  }

  function insertEmoji(unicode: string) {
    editorInstance.update(() => {
      const selection = ensureSelection();
      if ($isRangeSelection(selection)) {
        selection.insertText(unicode);
      } else {
        $insertNodes([$createTextNode(unicode)]);
      }
    });
    editorInstance.focus();
    setShowEmoji(false);
  }

  onMount(() => {
    editorInstance = createEditor({
      namespace: "pagoda-mini",
      nodes: [LinkNode],
      theme: {
        paragraph: "editor-paragraph",
        text: {
          bold: "editor-bold",
          italic: "editor-italic",
          underline: "editor-underline",
          highlight: "editor-spoiler",
        },
        link: "editor-link",
      },
      onError: (error: Error) => console.error(error),
    });

    editorInstance.setRootElement(editorRef);
    const cleanupRichText = registerRichText(editorInstance);
    const cleanupLink = editorInstance.registerCommand(TOGGLE_LINK_COMMAND, (payload) => {
      $toggleLink(typeof payload === "string" ? payload : null);
      return true;
    }, COMMAND_PRIORITY_LOW);
    const cleanupHistory = registerHistory(editorInstance, createEmptyHistoryState(), 300);

    if (props.initialHtml) {
      editorInstance.update(() => {
        const parser = new DOMParser();
        const dom = parser.parseFromString(props.initialHtml!, "text/html");
        const nodes = $generateNodesFromDOM(editorInstance, dom);
        const root = $getRoot();
        root.clear();
        root.append(...nodes);
      });
    }

    editorInstance.registerUpdateListener(({ editorState }) => {
      editorState.read(() => {
        props.onHtml($generateHtmlFromNodes(editorInstance));
        const selection = $getSelection();
        if ($isRangeSelection(selection)) {
          setBold(selection.hasFormat("bold"));
          setItalic(selection.hasFormat("italic"));
          setUnderline(selection.hasFormat("underline"));
          setSpoiler(selection.hasFormat("highlight"));

          const node = selection.anchor.getNode();
          const parent = node.getParent();
          setLink(parent !== null && $isLinkNode(parent));
        }
      });
    });

    function handleClickOutside(event: MouseEvent) {
      const target = event.target as HTMLElement;
      if (showEmoji() && !containerRef.contains(target) && !target.closest(".editor-emoji-popup")) {
        setShowEmoji(false);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);

    onCleanup(() => {
      cleanupRichText();
      cleanupLink();
      cleanupHistory();
      document.removeEventListener("mousedown", handleClickOutside);
      editorInstance.setRootElement(null);
    });
  });

  const s = "18";
  const w = "2";
  const wActive = "3";

  return (
    <div class="editor-container editor-mini" ref={containerRef}>
      <div class="editor-toolbar">
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": bold() }} onMouseDown={(e) => e.preventDefault()} onClick={() => formatText("bold")} title="Bold">
          <Show when={bold()} fallback={<IconBoldOff size={s} stroke={w} />}><IconBold size={s} stroke={w} /></Show>
        </button>
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": italic() }} onMouseDown={(e) => e.preventDefault()} onClick={() => formatText("italic")} title="Italic">
          <IconItalic size={s} stroke={italic() ? wActive : w} />
        </button>
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": underline() }} onMouseDown={(e) => e.preventDefault()} onClick={() => formatText("underline")} title="Underline">
          <IconUnderline size={s} stroke={underline() ? wActive : w} />
        </button>
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": spoiler() }} onMouseDown={(e) => e.preventDefault()} onClick={() => formatText("highlight")} title="Spoiler">
          <IconEyeOff size={s} stroke={spoiler() ? wActive : w} />
        </button>
        <span class="editor-toolbar-divider" />
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": link() }} onMouseDown={(e) => e.preventDefault()} onClick={openLinkInput} title="Link">
          <Show when={link()} fallback={<IconLinkOff size={s} stroke={w} />}><IconLink size={s} stroke={w} /></Show>
        </button>
        <span class="editor-toolbar-divider" />
        <button type="button" ref={emojiBtnRef} class="editor-toolbar-btn" onClick={toggleEmojiPicker} title="Emoji">
          <IconMoodSmile size={s} stroke={w} />
        </button>
      </div>
      <Show when={showLinkInput()}>
        <div class="editor-link-bar">
          <Show when={!hasSelection()}>
            <input
              type="text"
              class="editor-link-input"
              placeholder="Link text"
              value={linkText()}
              onInput={(e) => setLinkText(e.currentTarget.value)}
              onKeyDown={(e) => { if (e.key === "Escape") cancelLink(); }}
            />
          </Show>
          <input
            ref={linkUrlRef}
            type="text"
            class="editor-link-input"
            placeholder="https://"
            value={linkUrl()}
            onInput={(e) => setLinkUrl(e.currentTarget.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") submitLink();
              if (e.key === "Escape") cancelLink();
            }}
          />
          <button type="button" class="editor-link-btn editor-link-apply" onClick={submitLink}>Apply</button>
          <button type="button" class="editor-link-btn" onClick={cancelLink}>Cancel</button>
        </div>
      </Show>
      <div ref={editorRef} class="editor-root" contentEditable />
      <Show when={showEmoji()}>
        <div class="editor-emoji-popup" style={{ top: `${emojiPos().top}px`, right: `${emojiPos().right}px` }}>
          <emoji-picker
            class="editor-emoji-picker"
            on:emoji-click={(e: CustomEvent) => insertEmoji(e.detail.unicode)}
          />
        </div>
      </Show>
    </div>
  );
}