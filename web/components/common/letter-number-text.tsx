type LetterNumberTextProps = {
  value: string;
  className?: string;
};

function tokenizeLetterNumber(value: string) {
  return value
    .trim()
    .split(/([\-_/\\.\s]+)/g)
    .filter(Boolean);
}

export function LetterNumberText({ value, className }: LetterNumberTextProps) {
  const tokens = tokenizeLetterNumber(value);

  return (
    <span
      className={[
        "inline-flex flex-row items-center justify-start gap-0",
        className,
      ]
        .filter(Boolean)
        .join(" ")}
      dir="ltr"
      style={{
        direction: "ltr",
        unicodeBidi: "isolate",
      }}
    >
      {tokens.map((token, index) => (
        <span
          key={`${token}-${index}`}
          dir="ltr"
          style={{
            direction: "ltr",
            unicodeBidi: "isolate",
          }}
        >
          {token}
        </span>
      ))}
    </span>
  );
}