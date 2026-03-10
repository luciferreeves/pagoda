import { Show } from "solid-js";

interface PaginationProps {
  page: number;
  totalPages: number;
  total: number;
  label: string;
  onPage: (page: number) => void;
}

export default function Pagination(props: PaginationProps) {
  return (
    <Show when={props.totalPages > 1}>
      <div class="council-pagination">
        <button
          type="button"
          class="council-page-btn"
          disabled={props.page <= 1}
          onClick={() => props.onPage(props.page - 1)}
        >
          Prev
        </button>
        <span class="council-page-info">
          Page {props.page} of {props.totalPages} ({props.total} {props.label})
        </span>
        <button
          type="button"
          class="council-page-btn"
          disabled={props.page >= props.totalPages}
          onClick={() => props.onPage(props.page + 1)}
        >
          Next
        </button>
      </div>
    </Show>
  );
}