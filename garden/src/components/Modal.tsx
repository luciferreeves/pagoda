import type { JSX } from "solid-js";

interface ModalProps {
  title: string;
  onClose: () => void;
  children: JSX.Element;
}

export default function Modal(props: ModalProps) {
  function handleBackdrop(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      props.onClose();
    }
  }

  return (
    <div class="modal-backdrop" onClick={handleBackdrop}>
      <div class="modal-content">
        <div class="modal-header">
          <span class="modal-title">{props.title}</span>
          <button type="button" class="modal-close" onClick={props.onClose}>&times;</button>
        </div>
        <div class="modal-body">
          {props.children}
        </div>
      </div>
    </div>
  );
}