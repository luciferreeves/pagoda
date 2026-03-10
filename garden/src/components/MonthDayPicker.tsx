import { createSignal, Show, For } from "solid-js";
import { useClickOutside } from "../utils/clickOutside";

interface MonthDayPickerProps {
  value: string;
  onChange: (value: string) => void;
}

const MONTHS = [
  "Jan", "Feb", "Mar", "Apr", "May", "Jun",
  "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
];

function daysInMonth(month: number) {
  return new Date(2024, month + 1, 0).getDate();
}

function parseMonthDay(value: string) {
  if (!value) return null;
  const parts = value.split("-");
  return { month: parseInt(parts[0]) - 1, day: parseInt(parts[1]) };
}

function formatMonthDay(month: number, day: number) {
  return `${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
}

export default function MonthDayPicker(props: MonthDayPickerProps) {
  const initial = parseMonthDay(props.value);
  const [open, setOpen] = createSignal(false);
  const [selectedMonth, setSelectedMonth] = createSignal(initial?.month ?? -1);
  const [selectedDay, setSelectedDay] = createSignal(initial?.day ?? -1);
  const [viewMonth, setViewMonth] = createSignal(initial?.month ?? 0);
  let containerRef: HTMLDivElement | undefined;
  let monthColRef: HTMLDivElement | undefined;

  function dayCount() {
    return daysInMonth(viewMonth());
  }

  function scrollToSelectedMonth() {
    requestAnimationFrame(() => {
      const el = monthColRef?.querySelector(".mdp-wheel-item-selected") as HTMLElement | null;
      if (el && monthColRef) {
        monthColRef.scrollTop = el.offsetTop - monthColRef.offsetHeight / 2 + el.offsetHeight / 2;
      }
    });
  }

  function toggleOpen() {
    const next = !open();
    setOpen(next);
    if (next) scrollToSelectedMonth();
  }

  function selectMonth(value: number) {
    setViewMonth(value);
    setSelectedMonth(value);
    const d = selectedDay();
    if (d > 0) {
      const max = daysInMonth(value);
      const clamped = d > max ? max : d;
      setSelectedDay(clamped);
      props.onChange(formatMonthDay(value, clamped));
    }
  }

  function selectDay(day: number) {
    setSelectedDay(day);
    setSelectedMonth(viewMonth());
    props.onChange(formatMonthDay(viewMonth(), day));
  }

  function isSelectedDay(day: number) {
    return selectedMonth() === viewMonth() && selectedDay() === day;
  }

  function clear() {
    setSelectedMonth(-1);
    setSelectedDay(-1);
    props.onChange("");
    setOpen(false);
  }

  function displayValue() {
    const m = selectedMonth();
    const d = selectedDay();
    if (m < 0 || d < 0) return "—";
    return `${MONTHS[m]} ${d}`;
  }

  useClickOutside(() => containerRef, setOpen);

  return (
    <div class="mdp" ref={containerRef}>
      <button type="button" class="mdp-trigger" onClick={toggleOpen}>
        {displayValue()}
      </button>
      <Show when={open()}>
        <div class="mdp-popup">
          <div class="mdp-body">
            <div class="mdp-month-col" ref={monthColRef}>
              <For each={MONTHS}>
                {(name: string, i: () => number) => (
                  <button
                    type="button"
                    class="mdp-wheel-item"
                    classList={{ "mdp-wheel-item-selected": viewMonth() === i() }}
                    onClick={() => selectMonth(i())}
                  >
                    {name}
                  </button>
                )}
              </For>
            </div>
            <div class="mdp-day-col">
              <div class="mdp-grid">
                <For each={Array.from({ length: dayCount() }, (_, i) => i + 1)}>
                  {(d: number) => (
                    <button
                      type="button"
                      class="mdp-day"
                      classList={{ "mdp-day-selected": isSelectedDay(d) }}
                      onClick={() => selectDay(d)}
                    >
                      {d}
                    </button>
                  )}
                </For>
              </div>
            </div>
          </div>
          <div class="mdp-footer">
            <button type="button" class="mdp-clear" onClick={clear}>Clear</button>
          </div>
        </div>
      </Show>
    </div>
  );
}