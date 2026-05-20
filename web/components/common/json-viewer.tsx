export function JsonViewer({ value }: { value: unknown }) {
  if (value === undefined || value === null) {
    return (
      <div className="rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground">
        null
      </div>
    );
  }

  return (
    <pre className="max-h-96 overflow-auto rounded-md border bg-muted/30 p-3 text-xs">
      {JSON.stringify(value, null, 2)}
    </pre>
  );
}