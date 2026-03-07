import { createSignal, onMount, onCleanup, Show, For } from "solid-js";
import { createEditor, $getSelection, $isRangeSelection, $isElementNode, $createTextNode, $createParagraphNode, $insertNodes, $getRoot, COMMAND_PRIORITY_LOW, COMMAND_PRIORITY_HIGH, KEY_TAB_COMMAND, INDENT_CONTENT_COMMAND, OUTDENT_CONTENT_COMMAND } from "lexical";
import type { LexicalEditor, TextFormatType } from "lexical";
import { registerRichText, HeadingNode, QuoteNode, $createHeadingNode, $createQuoteNode, $isHeadingNode, $isQuoteNode } from "@lexical/rich-text";
import type { HeadingTagType } from "@lexical/rich-text";
import { $generateHtmlFromNodes } from "@lexical/html";
import { LinkNode, $createLinkNode, $isLinkNode, TOGGLE_LINK_COMMAND, $toggleLink } from "@lexical/link";
import { ListNode, ListItemNode, registerList, INSERT_UNORDERED_LIST_COMMAND, INSERT_ORDERED_LIST_COMMAND, REMOVE_LIST_COMMAND } from "@lexical/list";
import { $getNearestNodeOfType } from "@lexical/utils";
import { $setBlocksType } from "@lexical/selection";
import { CodeNode, CodeHighlightNode, registerCodeHighlighting, $createCodeNode, $isCodeNode } from "@lexical/code";
import { registerHistory, createEmptyHistoryState } from "@lexical/history";
import {
  IconBold, IconBoldOff,
  IconItalic,
  IconUnderline,
  IconStrikethrough,
  IconCode, IconCodeOff,
  IconLink, IconLinkOff,
  IconList, IconListNumbers,
  IconBlockquote,
  IconSourceCode,
  IconH1, IconH2, IconH3,
  IconHeading, IconHeadingOff,
  IconEyeOff,
  IconMoodSmile,
  IconPaperclip,
  IconPhoto,
  IconFile,
  IconFileText,
  IconFileSpreadsheet,
  IconPresentation,
  IconFileCode,
  IconMusic,
  IconVideo,
  IconFileZip,
  IconFileTypePdf,
  IconTypography,
  IconDatabase,
  IconX,
} from "@tabler/icons-solidjs";
import "emoji-picker-element";
import { api } from "../api";
import { auth } from "../store/auth";
import { API_URL } from "../config";

interface AttachmentResult {
  ref: string;
  file_name: string;
  url: string;
  file_size: number;
  content_type: string;
}

interface EditorProps {
  onHtml: (html: string) => void;
  onAttachment?: (ref: string) => void;
}

export default function Editor(props: EditorProps) {
  let editorRef!: HTMLDivElement;
  let containerRef!: HTMLDivElement;
  let linkUrlRef!: HTMLInputElement;
  let editorInstance: LexicalEditor;

  const [bold, setBold] = createSignal(false);
  const [italic, setItalic] = createSignal(false);
  const [underline, setUnderline] = createSignal(false);
  const [strikethrough, setStrikethrough] = createSignal(false);
  const [code, setCode] = createSignal(false);
  const [spoiler, setSpoiler] = createSignal(false);
  const [link, setLink] = createSignal(false);
  const [bulletList, setBulletList] = createSignal(false);
  const [orderedList, setOrderedList] = createSignal(false);
  const [blockQuote, setBlockQuote] = createSignal(false);
  const [codeBlock, setCodeBlock] = createSignal(false);
  const [heading, setHeading] = createSignal<string>("");
  const [showLinkInput, setShowLinkInput] = createSignal(false);
  const [linkUrl, setLinkUrl] = createSignal("");
  const [linkText, setLinkText] = createSignal("");
  const [hasSelection, setHasSelection] = createSignal(false);
  const [showHeadings, setShowHeadings] = createSignal(false);
  const [showEmoji, setShowEmoji] = createSignal(false);
  const [emojiPos, setEmojiPos] = createSignal({ top: 0, right: 0 });
  type PendingUpload = { id: number; name: string; type: string; preview?: string; progress: number };
  let uploadId = 0;
  const [pendingUploads, setPendingUploads] = createSignal<PendingUpload[]>([]);
  const [attachments, setAttachments] = createSignal<AttachmentResult[]>([]);
  const [maxFileSize, setMaxFileSize] = createSignal(33554432);
  const [maxAttachments, setMaxAttachments] = createSignal(8);
  const [uploadError, setUploadError] = createSignal("");
  let fileInputRef!: HTMLInputElement;

  api<{ max_file_size: number; max_attachments: number }>("/config").then((res) => {
    if (res.ok) {
      setMaxFileSize(res.data.max_file_size);
      setMaxAttachments(res.data.max_attachments);
    }
  });

  function handleFileSelect() {
    const file = fileInputRef.files?.[0];
    if (!file) return;
    fileInputRef.value = "";

    setUploadError("");
    if (file.size > maxFileSize()) {
      setUploadError(`File exceeds the maximum size of ${Math.floor(maxFileSize() / 1048576)} MB.`);
      return;
    }
    if (attachments().length + pendingUploads().length >= maxAttachments()) {
      setUploadError(`Maximum of ${maxAttachments()} attachments allowed.`);
      return;
    }

    const token = auth.token();
    if (!token) return;

    const id = ++uploadId;
    let preview: string | undefined;
    if (file.type.startsWith("image/") || file.type.startsWith("video/")) {
      preview = URL.createObjectURL(file);
    }

    setPendingUploads((prev: PendingUpload[]) => [...prev, { id, name: file.name, type: file.type, preview, progress: 0 }]);

    const xhr = new XMLHttpRequest();
    const form = new FormData();
    form.append("file", file);

    xhr.upload.addEventListener("progress", (e) => {
      if (e.lengthComputable) {
        const pct = Math.round((e.loaded / e.total) * 100);
        setPendingUploads((prev: PendingUpload[]) => prev.map((p: PendingUpload) => p.id === id ? { ...p, progress: pct } : p));
      }
    });

    xhr.addEventListener("load", () => {
      const pending = pendingUploads().find((p: PendingUpload) => p.id === id);
      if (pending?.preview) URL.revokeObjectURL(pending.preview);
      setPendingUploads((prev: PendingUpload[]) => prev.filter((p: PendingUpload) => p.id !== id));

      if (xhr.status < 200 || xhr.status >= 300) return;

      const attachment: AttachmentResult = JSON.parse(xhr.responseText);
      setAttachments((prev: AttachmentResult[]) => [...prev, attachment]);
      props.onAttachment?.(attachment.ref);
    });

    xhr.addEventListener("error", () => {
      const errPending = pendingUploads().find((p: PendingUpload) => p.id === id);
      if (errPending?.preview) URL.revokeObjectURL(errPending.preview);
      setPendingUploads((prev: PendingUpload[]) => prev.filter((p: PendingUpload) => p.id !== id));
    });

    xhr.open("POST", `${API_URL}/letters/attachments`);
    xhr.setRequestHeader("Authorization", `Bearer ${token}`);
    xhr.send(form);
  }

  function removeAttachment(ref: string) {
    setAttachments((prev: AttachmentResult[]) => prev.filter((a: AttachmentResult) => a.ref !== ref));
  }

  function formatFileSize(bytes: number): string {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1048576) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / 1048576).toFixed(1)} MB`;
  }

  function fileExtension(name: string): string {
    const dot = name.lastIndexOf(".");
    return dot >= 0 ? name.slice(dot + 1).toUpperCase() : "";
  }

  const archiveMimes = ["application/zip", "application/x-rar-compressed", "application/gzip", "application/x-7z-compressed", "application/x-tar", "application/x-bzip2", "application/x-xz", "application/x-compress"];
  const codeMimes = ["application/javascript", "application/json", "application/xml", "application/x-httpd-php", "application/x-sh", "application/x-python", "application/typescript"];
  const docMimes = ["application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.oasis.opendocument.text", "application/rtf", "application/epub+zip"];
  const sheetMimes = ["application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.oasis.opendocument.spreadsheet", "text/csv"];
  const slideMimes = ["application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation", "application/vnd.oasis.opendocument.presentation"];
  const dbMimes = ["application/x-sqlite3", "application/vnd.ms-access"];

  function fileIcon(contentType: string) {
    const s = "24";
    const k = "1.5";
    if (contentType.startsWith("image/")) return <IconPhoto size={s} stroke={k} />;
    if (contentType.startsWith("audio/")) return <IconMusic size={s} stroke={k} />;
    if (contentType.startsWith("video/")) return <IconVideo size={s} stroke={k} />;
    if (contentType.startsWith("font/") || contentType === "application/font-sfnt" || contentType === "application/vnd.ms-fontobject" || contentType === "application/font-woff" || contentType === "application/font-woff2") return <IconTypography size={s} stroke={k} />;
    if (contentType === "application/pdf") return <IconFileTypePdf size={s} stroke={k} />;
    if (archiveMimes.includes(contentType)) return <IconFileZip size={s} stroke={k} />;
    if (dbMimes.includes(contentType)) return <IconDatabase size={s} stroke={k} />;
    if (docMimes.includes(contentType)) return <IconFileText size={s} stroke={k} />;
    if (sheetMimes.includes(contentType)) return <IconFileSpreadsheet size={s} stroke={k} />;
    if (slideMimes.includes(contentType)) return <IconPresentation size={s} stroke={k} />;
    if (codeMimes.includes(contentType) || contentType.startsWith("text/x-")) return <IconFileCode size={s} stroke={k} />;
    if (contentType.startsWith("text/")) return <IconFileText size={s} stroke={k} />;
    return <IconFile size={s} stroke={k} />;
  }

  function closeAllPopups() {
    setShowHeadings(false);
    setShowEmoji(false);
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
    closeAllPopups();
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

  function ensureSelection() {
    let selection = $getSelection();
    if (!$isRangeSelection(selection)) {
      const root = $getRoot();
      const firstChild = root.getFirstChild();
      if ($isElementNode(firstChild)) {
        selection = firstChild.select(0, 0);
      } else {
        const p = $createParagraphNode();
        root.append(p);
        selection = p.select(0, 0);
      }
    }
    return selection;
  }

  function toggleBulletList() {
    editorInstance.update(() => { ensureSelection(); });
    if (bulletList()) {
      editorInstance.dispatchCommand(REMOVE_LIST_COMMAND, undefined);
    } else {
      editorInstance.dispatchCommand(INSERT_UNORDERED_LIST_COMMAND, undefined);
    }
    editorInstance.focus();
  }

  function toggleOrderedList() {
    editorInstance.update(() => { ensureSelection(); });
    if (orderedList()) {
      editorInstance.dispatchCommand(REMOVE_LIST_COMMAND, undefined);
    } else {
      editorInstance.dispatchCommand(INSERT_ORDERED_LIST_COMMAND, undefined);
    }
    editorInstance.focus();
  }

  function toggleBlockQuote() {
    editorInstance.update(() => {
      const selection = ensureSelection();
      if ($isRangeSelection(selection)) {
        if (blockQuote()) {
          $setBlocksType(selection, () => $createParagraphNode());
        } else {
          $setBlocksType(selection, () => $createQuoteNode());
        }
      }
    });
    editorInstance.focus();
  }

  function toggleCodeBlock() {
    editorInstance.update(() => {
      const selection = ensureSelection();
      if ($isRangeSelection(selection)) {
        if (codeBlock()) {
          $setBlocksType(selection, () => $createParagraphNode());
        } else {
          $setBlocksType(selection, () => $createCodeNode());
        }
      }
    });
    editorInstance.focus();
  }

  function setHeadingLevel(tag: HeadingTagType | null) {
    editorInstance.update(() => {
      const selection = ensureSelection();
      if ($isRangeSelection(selection)) {
        if (tag) {
          $setBlocksType(selection, () => $createHeadingNode(tag));
        } else {
          $setBlocksType(selection, () => $createParagraphNode());
        }
      }
    });
    editorInstance.focus();
    setShowHeadings(false);
  }

  function toggleHeadingsMenu() {
    setShowEmoji(false);
    setShowHeadings(!showHeadings());
  }

  let emojiBtnRef!: HTMLButtonElement;

  function toggleEmojiPicker() {
    setShowHeadings(false);
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


  function formatSpoiler() {
    formatText("highlight");
  }

  onMount(() => {
    editorInstance = createEditor({
      namespace: "pagoda",
      nodes: [HeadingNode, QuoteNode, LinkNode, ListNode, ListItemNode, CodeNode, CodeHighlightNode],
      theme: {
        paragraph: "editor-paragraph",
        text: {
          bold: "editor-bold",
          italic: "editor-italic",
          underline: "editor-underline",
          strikethrough: "editor-strikethrough",
          code: "editor-code",
          highlight: "editor-spoiler",
        },
        link: "editor-link",
        list: {
          ul: "editor-ul",
          ol: "editor-ol",
          listitem: "editor-li",
          nested: { listitem: "editor-li-nested" },
        },
        quote: "editor-quote",
        code: "editor-codeblock",
        codeHighlight: {},
        heading: {
          h1: "editor-h1",
          h2: "editor-h2",
          h3: "editor-h3",
        },
      },
      onError: (error: Error) => console.error(error),
    });

    editorInstance.setRootElement(editorRef);
    const cleanupRichText = registerRichText(editorInstance);
    const cleanupList = registerList(editorInstance);
    const cleanupLink = editorInstance.registerCommand(TOGGLE_LINK_COMMAND, (payload) => {
      $toggleLink(typeof payload === "string" ? payload : null);
      return true;
    }, COMMAND_PRIORITY_LOW);
    const cleanupCode = registerCodeHighlighting(editorInstance, undefined);
    const cleanupHistory = registerHistory(editorInstance, createEmptyHistoryState(), 300);

    const cleanupTab = editorInstance.registerCommand(KEY_TAB_COMMAND, (event) => {
      event.preventDefault();
      if (event.shiftKey) {
        editorInstance.dispatchCommand(OUTDENT_CONTENT_COMMAND, undefined);
      } else {
        editorInstance.dispatchCommand(INDENT_CONTENT_COMMAND, undefined);
      }
      return true;
    }, COMMAND_PRIORITY_HIGH);

    editorInstance.registerUpdateListener(({ editorState }) => {
      editorState.read(() => {
        props.onHtml($generateHtmlFromNodes(editorInstance));
        const selection = $getSelection();
        if ($isRangeSelection(selection)) {
          setBold(selection.hasFormat("bold"));
          setItalic(selection.hasFormat("italic"));
          setUnderline(selection.hasFormat("underline"));
          setStrikethrough(selection.hasFormat("strikethrough"));
          setCode(selection.hasFormat("code"));
          setSpoiler(selection.hasFormat("highlight"));

          const node = selection.anchor.getNode();
          const parent = node.getParent();
          setLink(parent !== null && $isLinkNode(parent));

          const listItem = $getNearestNodeOfType(node, ListItemNode);
          if (listItem) {
            const listNode = listItem.getParent();
            if (listNode instanceof ListNode) {
              const tag = listNode.getListType();
              setBulletList(tag === "bullet");
              setOrderedList(tag === "number");
            } else {
              setBulletList(false);
              setOrderedList(false);
            }
          } else {
            setBulletList(false);
            setOrderedList(false);
          }

          const anchorNode = selection.anchor.getNode();
          const topLevel = anchorNode.getKey() === "root" ? anchorNode : anchorNode.getTopLevelElementOrThrow();
          setBlockQuote($isQuoteNode(topLevel));
          setCodeBlock($isCodeNode(topLevel));
          if ($isHeadingNode(topLevel)) {
            setHeading(topLevel.getTag());
          } else {
            setHeading("");
          }
        }
      });
    });

    function handleClickOutside(e: MouseEvent) {
      const target = e.target as HTMLElement;
      if (!containerRef.contains(target)) {
        if (showEmoji() && !target.closest(".editor-emoji-popup")) {
          setShowEmoji(false);
        }
        if (showHeadings()) {
          setShowHeadings(false);
        }
      } else {
        if (showHeadings() && !target.closest(".editor-toolbar-dropdown-wrap:has(.editor-dropdown)")) {
          setShowHeadings(false);
        }
        if (showEmoji() && !target.closest(".editor-emoji-popup")) {
          setShowEmoji(false);
        }
      }
    }

    document.addEventListener("mousedown", handleClickOutside);

    onCleanup(() => {
      cleanupRichText();
      cleanupList();
      cleanupLink();
      cleanupCode();
      cleanupHistory();
      cleanupTab();
      document.removeEventListener("mousedown", handleClickOutside);
      editorInstance.setRootElement(null);
    });
  });

  const s = "18";
  const w = "2";
  const wActive = "3";

  return (
    <div class="editor-container" ref={containerRef}>
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
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": strikethrough() }} onMouseDown={(e) => e.preventDefault()} onClick={() => formatText("strikethrough")} title="Strikethrough">
          <IconStrikethrough size={s} stroke={strikethrough() ? wActive : w} />
        </button>
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": spoiler() }} onMouseDown={(e) => e.preventDefault()} onClick={formatSpoiler} title="Spoiler">
          <IconEyeOff size={s} stroke={spoiler() ? wActive : w} />
        </button>
        <span class="editor-toolbar-divider" />
        <div class="editor-toolbar-dropdown-wrap">
          <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": heading() !== "" }} onMouseDown={(e) => e.preventDefault()} onClick={toggleHeadingsMenu} title="Heading">
            <Show when={heading() !== ""} fallback={<IconHeadingOff size={s} stroke={w} />}><IconHeading size={s} stroke={w} /></Show>
          </button>
          <Show when={showHeadings()}>
            <div class="editor-dropdown">
              <button type="button" class="editor-dropdown-item" classList={{ "editor-dropdown-active": heading() === "h1" }} onMouseDown={(e) => e.preventDefault()} onClick={() => setHeadingLevel("h1")}><IconH1 size="16" stroke={w} /> Heading 1</button>
              <button type="button" class="editor-dropdown-item" classList={{ "editor-dropdown-active": heading() === "h2" }} onMouseDown={(e) => e.preventDefault()} onClick={() => setHeadingLevel("h2")}><IconH2 size="16" stroke={w} /> Heading 2</button>
              <button type="button" class="editor-dropdown-item" classList={{ "editor-dropdown-active": heading() === "h3" }} onMouseDown={(e) => e.preventDefault()} onClick={() => setHeadingLevel("h3")}><IconH3 size="16" stroke={w} /> Heading 3</button>
              <button type="button" class="editor-dropdown-item" onMouseDown={(e) => e.preventDefault()} onClick={() => setHeadingLevel(null)}>Normal text</button>
            </div>
          </Show>
        </div>
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": blockQuote() }} onMouseDown={(e) => e.preventDefault()} onClick={toggleBlockQuote} title="Block Quote">
          <IconBlockquote size={s} stroke={blockQuote() ? wActive : w} />
        </button>
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": codeBlock() }} onMouseDown={(e) => e.preventDefault()} onClick={toggleCodeBlock} title="Code Block">
          <IconSourceCode size={s} stroke={codeBlock() ? wActive : w} />
        </button>
        <span class="editor-toolbar-divider" />
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": code() }} onMouseDown={(e) => e.preventDefault()} onClick={() => formatText("code")} title="Inline Code">
          <Show when={code()} fallback={<IconCodeOff size={s} stroke={w} />}><IconCode size={s} stroke={w} /></Show>
        </button>
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": link() }} onMouseDown={(e) => e.preventDefault()} onClick={openLinkInput} title="Link">
          <Show when={link()} fallback={<IconLinkOff size={s} stroke={w} />}><IconLink size={s} stroke={w} /></Show>
        </button>
        <span class="editor-toolbar-divider" />
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": bulletList() }} onMouseDown={(e) => e.preventDefault()} onClick={toggleBulletList} title="Bullet List">
          <IconList size={s} stroke={bulletList() ? wActive : w} />
        </button>
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-active": orderedList() }} onMouseDown={(e) => e.preventDefault()} onClick={toggleOrderedList} title="Numbered List">
          <IconListNumbers size={s} stroke={orderedList() ? wActive : w} />
        </button>
        <span class="editor-toolbar-divider" />
        <button type="button" ref={emojiBtnRef} class="editor-toolbar-btn" onClick={toggleEmojiPicker} title="Emoji">
          <IconMoodSmile size={s} stroke={w} />
        </button>
        <button type="button" class="editor-toolbar-btn" classList={{ "editor-toolbar-disabled": attachments().length + pendingUploads().length >= maxAttachments() }} onMouseDown={(e) => e.preventDefault()} onClick={() => fileInputRef.click()} title="Attach File" disabled={attachments().length + pendingUploads().length >= maxAttachments()}>
          <IconPaperclip size={s} stroke={w} />
        </button>
        <input ref={fileInputRef} type="file" class="editor-file-input" onChange={handleFileSelect} />
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
              onKeyDown={(e) => {
                if (e.key === "Escape") cancelLink();
              }}
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
      <Show when={uploadError()}>
        <div class="editor-upload-error">{uploadError()}</div>
      </Show>
      <Show when={attachments().length > 0 || pendingUploads().length > 0}>
        <div class="editor-attachments">
          <For each={attachments()}>
            {(attachment) => (
              <div class="editor-attachment" title={attachment.file_name}>
                <Show when={attachment.content_type.startsWith("image/")} fallback={
                  <Show when={attachment.content_type.startsWith("video/")} fallback={
                    <div class="editor-attachment-tile">
                      {fileIcon(attachment.content_type)}
                      <span class="editor-attachment-ext">{fileExtension(attachment.file_name)}</span>
                    </div>
                  }>
                    <video src={attachment.url} class="editor-attachment-thumb" preload="metadata" />
                    <div class="editor-attachment-badge"><IconVideo size="12" stroke="2" /></div>
                  </Show>
                }>
                  <img src={attachment.url} alt={attachment.file_name} class="editor-attachment-thumb" />
                  <div class="editor-attachment-badge"><IconPhoto size="12" stroke="2" /></div>
                </Show>
                <span class="editor-attachment-size">{formatFileSize(attachment.file_size)}</span>
                <button type="button" class="editor-attachment-remove" onClick={() => removeAttachment(attachment.ref)} title="Remove">
                  <IconX size="12" stroke="2" />
                </button>
              </div>
            )}
          </For>
          <For each={pendingUploads()}>
            {(pending) => (
              <div class="editor-attachment editor-attachment-uploading" title={pending.name}>
                <Show when={pending.type.startsWith("image/") && pending.preview} fallback={
                  <Show when={pending.type.startsWith("video/") && pending.preview} fallback={
                    <div class="editor-attachment-tile">
                      {fileIcon(pending.type)}
                      <span class="editor-attachment-ext">{fileExtension(pending.name)}</span>
                    </div>
                  }>
                    <video src={pending.preview} class="editor-attachment-thumb" preload="metadata" />
                    <div class="editor-attachment-badge"><IconVideo size="12" stroke="2" /></div>
                  </Show>
                }>
                  <img src={pending.preview} alt={pending.name} class="editor-attachment-thumb" />
                  <div class="editor-attachment-badge"><IconPhoto size="12" stroke="2" /></div>
                </Show>
                <div class="editor-attachment-pending-bar">
                  <div class="editor-attachment-pending-fill" style={{ width: `${pending.progress}%` }} />
                </div>
              </div>
            )}
          </For>
        </div>
      </Show>
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