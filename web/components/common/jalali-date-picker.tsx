"use client";

import DatePicker, { DateObject } from "react-multi-date-picker";
import persian from "react-date-object/calendars/persian";
import persianFa from "react-date-object/locales/persian_fa";

type JalaliDatePickerProps = {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  required?: boolean;
  id?: string;
};

export function JalaliDatePicker({
  value,
  onChange,
  placeholder = "۱۴۰۵/۰۳/۰۱",
  required,
  id,
}: JalaliDatePickerProps) {
  const pickerValue = value
    ? new DateObject({
        date: value,
        calendar: persian,
        locale: persianFa,
        format: "YYYY/MM/DD",
      })
    : undefined;

  return (
    <DatePicker
      id={id}
      value={pickerValue}
      onChange={(date) => {
        if (!date) {
          onChange("");
          return;
        }

        const selectedDate = Array.isArray(date) ? date[0] : date;
        onChange(selectedDate.format("YYYY/MM/DD"));
      }}
      calendar={persian}
      locale={persianFa}
      format="YYYY/MM/DD"
      calendarPosition="bottom-right"
      containerClassName="w-full"
      inputClass="rmdp-input"
      placeholder={placeholder}
      editable={false}
      required={required}
      weekDays={["ش", "ی", "د", "س", "چ", "پ", "ج"]}
      months={[
        "فروردین",
        "اردیبهشت",
        "خرداد",
        "تیر",
        "مرداد",
        "شهریور",
        "مهر",
        "آبان",
        "آذر",
        "دی",
        "بهمن",
        "اسفند",
      ]}
    />
  );
}