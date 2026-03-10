import { createSignal, Show, For } from "solid-js";
import { useClickOutside } from "../utils/clickOutside";

interface DatePickerProps {
  value: string;
  onChange: (value: string) => void;
}

const MONTHS = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];

const WEEKDAYS = ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"];

function daysInMonth(year: number, month: number) {
  return new Date(year, month + 1, 0).getDate();
}

function firstDayOfMonth(year: number, month: number) {
  return new Date(year, month, 1).getDay();
}

function parseDate(value: string) {
  if (!value) return null;
  const parts = value.split("-");
  return { year: parseInt(parts[0]), month: parseInt(parts[1]) - 1, day: parseInt(parts[2]) };
}

function formatDate(year: number, month: number, day: number) {
  return `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
}

export default function DatePicker(props: DatePickerProps) {
  const parsed = () => parseDate(props.value);
  const today = new Date();
  let triggerRef: HTMLButtonElement | undefined;
  let containerRef: HTMLDivElement | undefined;

  const [open, setOpen] = createSignal(false);
  const [dropdownStyle, setDropdownStyle] = createSignal<Record<string, string>>({});
  const [viewYear, setViewYear] = createSignal(parsed()?.year ?? today.getFullYear());
  const [viewMonth, setViewMonth] = createSignal(parsed()?.month ?? today.getMonth());

  function toggleOpen() {
    if (!open() && triggerRef) {
      const rect = triggerRef.getBoundingClientRect();
      setDropdownStyle({ top: `${rect.bottom + 2}px`, left: `${rect.left}px` });
    }
    setOpen(!open());
  }

  useClickOutside(() => containerRef, setOpen);

  function prevMonth() {
    if (viewMonth() === 0) {
      setViewMonth(11);
      setViewYear(viewYear() - 1);
    } else {
      setViewMonth(viewMonth() - 1);
    }
  }

  function nextMonth() {
    if (viewMonth() === 11) {
      setViewMonth(0);
      setViewYear(viewYear() + 1);
    } else {
      setViewMonth(viewMonth() + 1);
    }
  }

  function selectDay(day: number) {
    props.onChange(formatDate(viewYear(), viewMonth(), day));
    setOpen(false);
  }

  function clear() {
    props.onChange("");
    setOpen(false);
  }

  function calendarDays() {
    const total = daysInMonth(viewYear(), viewMonth());
    const offset = firstDayOfMonth(viewYear(), viewMonth());
    const days: (number | null)[] = [];
    for (let blank = 0; blank < offset; blank++) days.push(null);
    for (let day = 1; day <= total; day++) days.push(day);
    return days;
  }

  function isSelected(day: number) {
    const selected = parsed();
    return selected && selected.year === viewYear() && selected.month === viewMonth() && selected.day === day;
  }

  function displayValue() {
    const selected = parsed();
    if (!selected) return "—";
    return `${MONTHS[selected.month]} ${selected.day}, ${selected.year}`;
  }

  return (
    <div class="datepicker" ref={containerRef}>
      <button type="button" class="datepicker-trigger" ref={triggerRef} onClick={toggleOpen}>
        {displayValue()}
      </button>
      <Show when={open()}>
        <div class="datepicker-dropdown datepicker-dropdown-fixed" style={dropdownStyle()}>
          <div class="datepicker-nav">
            <button type="button" class="datepicker-nav-btn" onClick={prevMonth}>&lsaquo;</button>
            <span class="datepicker-nav-title">{MONTHS[viewMonth()]} {viewYear()}</span>
            <button type="button" class="datepicker-nav-btn" onClick={nextMonth}>&rsaquo;</button>
          </div>
          <div class="datepicker-weekdays">
            <For each={WEEKDAYS}>
              {(day: string) => <span class="datepicker-weekday">{day}</span>}
            </For>
          </div>
          <div class="datepicker-grid">
            <For each={calendarDays()}>
              {(day: number | null) => (
                <Show when={day !== null} fallback={<span class="datepicker-blank" />}>
                  <button
                    type="button"
                    class="datepicker-day"
                    classList={{ "datepicker-day-selected": isSelected(day!) ?? undefined }}
                    onClick={() => selectDay(day!)}
                  >
                    {day}
                  </button>
                </Show>
              )}
            </For>
          </div>
          <div class="datepicker-footer">
            <button type="button" class="datepicker-clear" onClick={clear}>Clear</button>
          </div>
        </div>
      </Show>
    </div>
  );
}