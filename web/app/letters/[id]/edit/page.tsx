"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { toast } from "sonner";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { AppShell } from "@/components/layout/app-shell";
import { LetterNumberText } from "@/components/common/letter-number-text";
import { useAuth } from "@/contexts/auth-context";
import {
  deleteLetterAttachment,
  getAttachmentDownloadURL,
  listLetterAttachments,
  uploadLetterAttachments,
} from "@/lib/api/attachments";
import { getLetter, updateLetter } from "@/lib/api/letters";
import { useI18n } from "@/lib/i18n/i18n-context";
import type { LetterAttachment } from "@/types/attachment";
import type { Letter } from "@/types/letter";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { JalaliDatePicker } from "@/components/common/jalali-date-picker";

export default function EditLetterPage() {
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const { token } = useAuth();
  const { t } = useI18n();

  const letterID = params.id;

  const [letter, setLetter] = useState<Letter | null>(null);

  const [displayLetterNumber, setDisplayLetterNumber] = useState("");
  const [title, setTitle] = useState("");
  const [letterDate, setLetterDate] = useState("");
  const [sender, setSender] = useState("");
  const [receiver, setReceiver] = useState("");
  const [description, setDescription] = useState("");

  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  const [attachments, setAttachments] = useState<LetterAttachment[]>([]);
  const [attachmentsLoading, setAttachmentsLoading] = useState(false);
  const [selectedAttachmentFiles, setSelectedAttachmentFiles] = useState<File[]>(
    []
  );
  const [uploadingAttachments, setUploadingAttachments] = useState(false);
  const [deletingAttachmentID, setDeletingAttachmentID] = useState<string | null>(
    null
  );

  const loadLetter = useCallback(async () => {
    if (!token || !letterID) return;

    setLoading(true);

    try {
      const result = await getLetter(token, letterID);

      setLetter(result);
      setDisplayLetterNumber(result.formatted_letter_number || "");
      setTitle(result.title || "");
      setLetterDate(result.letter_date_jalali || "");
      setSender(result.sender || "");
      setReceiver(result.receiver || "");
      setDescription(result.description || "");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load letter");
    } finally {
      setLoading(false);
    }
  }, [token, letterID]);

  const loadAttachments = useCallback(async () => {
    if (!token || !letterID) return;

    setAttachmentsLoading(true);

    try {
      const result = await listLetterAttachments(token, letterID);
      setAttachments(result);
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to load attachments"
      );
    } finally {
      setAttachmentsLoading(false);
    }
  }, [token, letterID]);

  useEffect(() => {
    const timeoutID = window.setTimeout(() => {
      void loadLetter();
    }, 0);

    return () => {
      window.clearTimeout(timeoutID);
    };
  }, [loadLetter]);

  useEffect(() => {
    const timeoutID = window.setTimeout(() => {
      void loadAttachments();
      setSelectedAttachmentFiles([]);
    }, 0);

    return () => {
      window.clearTimeout(timeoutID);
    };
  }, [loadAttachments]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!token || !letterID) return;

    setSaving(true);

    try {
      const updated = await updateLetter(token, letterID, {
        display_letter_number: displayLetterNumber.trim(),
        title,
        letter_date: letterDate,
        sender,
        receiver,
        description: description.trim() || null,
      });

      toast.success("نامه با موفقیت ویرایش شد");
      setLetter(updated);

      router.push("/letters");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update letter");
    } finally {
      setSaving(false);
    }
  }

  async function handleOpenAttachment(attachment: LetterAttachment) {
    if (!token || !letterID) return;

    try {
      const result = await getAttachmentDownloadURL(
        token,
        letterID,
        attachment.id
      );

      window.open(result.url, "_blank", "noopener,noreferrer");
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to open attachment"
      );
    }
  }

  async function handleUploadAttachments() {
    if (!token || !letterID || selectedAttachmentFiles.length === 0) return;

    setUploadingAttachments(true);

    try {
      await uploadLetterAttachments(token, letterID, selectedAttachmentFiles);

      toast.success("فایل‌های پیوست با موفقیت آپلود شدند");
      setSelectedAttachmentFiles([]);

      await loadAttachments();
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Attachment upload failed"
      );
    } finally {
      setUploadingAttachments(false);
    }
  }

  async function handleDeleteAttachment(attachment: LetterAttachment) {
    if (!token || !letterID) return;

    setDeletingAttachmentID(attachment.id);

    try {
      await deleteLetterAttachment(token, letterID, attachment.id);

      setAttachments((current) =>
        current.filter((item) => item.id !== attachment.id)
      );

      toast.success("فایل پیوست حذف شد");
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to delete attachment"
      );
    } finally {
      setDeletingAttachmentID(null);
    }
  }

  return (
    <ProtectedRoute allowedRoles={["superuser", "editor"]}>
      <AppShell>
        <div className="max-w-3xl space-y-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              ویرایش نامه
            </h1>
            <p className="text-muted-foreground">
              اطلاعات نامه را ویرایش و ذخیره کنید.
            </p>
          </div>

          {letter ? (
            <Card>
              <CardHeader>
                <CardTitle>شماره فعلی نامه</CardTitle>
              </CardHeader>
              <CardContent>
                <LetterNumberText
                  value={letter.formatted_letter_number}
                  className="text-3xl font-bold"
                />
              </CardContent>
            </Card>
          ) : null}

          <Card>
            <CardHeader>
              <CardTitle>{t.letterInformation}</CardTitle>
            </CardHeader>

            <CardContent>
              {loading ? (
                <div className="text-sm text-muted-foreground">
                  {t.commonLoading}
                </div>
              ) : (
                <form id="edit-letter-form" onSubmit={handleSubmit} className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="display_letter_number">
                      {t.displayLetterNumber}
                    </Label>
                    <Input
                      id="display_letter_number"
                      value={displayLetterNumber}
                      onChange={(event) =>
                        setDisplayLetterNumber(event.target.value)
                      }
                      required
                      className="text-left"
                      style={{
                        direction: "ltr",
                        unicodeBidi: "plaintext",
                      }}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="title">{t.letterTitle}</Label>
                    <Input
                      id="title"
                      value={title}
                      onChange={(event) => setTitle(event.target.value)}
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="letter_date">{t.letterDate}</Label>
                    <JalaliDatePicker
                      id="letter_date"
                      value={letterDate}
                      onChange={setLetterDate}
                      required
                    />
                  </div>

                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <Label htmlFor="sender">{t.sender}</Label>
                      <Input
                        id="sender"
                        value={sender}
                        onChange={(event) => setSender(event.target.value)}
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="receiver">{t.receiver}</Label>
                      <Input
                        id="receiver"
                        value={receiver}
                        onChange={(event) => setReceiver(event.target.value)}
                        required
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="description">{t.description}</Label>
                    <Textarea
                      id="description"
                      value={description}
                      onChange={(event) => setDescription(event.target.value)}
                    />
                  </div>
                </form>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>پیوست‌ها / اسکن نامه</CardTitle>
            </CardHeader>

            <CardContent className="space-y-4">
              <div className="space-y-2 rounded-md border bg-muted/30 p-3">
                <Label htmlFor="attachments">افزودن فایل جدید</Label>

                <Input
                  id="attachments"
                  type="file"
                  multiple
                  accept=".pdf,image/jpeg,image/png"
                  onChange={(event) => {
                    setSelectedAttachmentFiles(
                      Array.from(event.target.files || [])
                    );
                  }}
                />

                <div className="flex flex-wrap items-center justify-between gap-2">
                  <p className="text-xs text-muted-foreground">
                    آپلود فایل اختیاری است. فرمت‌های مجاز: PDF, JPG, PNG. برای
                    جایگزینی فایل، فایل قبلی را حذف و فایل جدید را آپلود کنید.
                  </p>

                  <Button
                    type="button"
                    size="sm"
                    disabled={
                      selectedAttachmentFiles.length === 0 ||
                      uploadingAttachments
                    }
                    onClick={() => void handleUploadAttachments()}
                  >
                    {uploadingAttachments ? t.commonLoading : "آپلود فایل"}
                  </Button>
                </div>

                {selectedAttachmentFiles.length > 0 ? (
                  <ul className="space-y-1 text-xs text-muted-foreground">
                    {selectedAttachmentFiles.map((file) => (
                      <li key={`${file.name}-${file.size}`}>
                        {file.name} - {formatFileSize(file.size)}
                      </li>
                    ))}
                  </ul>
                ) : null}
              </div>

              <div className="space-y-3">
                <div className="font-semibold">فایل‌های فعلی</div>

                {attachmentsLoading ? (
                  <div className="text-sm text-muted-foreground">
                    {t.commonLoading}
                  </div>
                ) : attachments.length === 0 ? (
                  <div className="rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground">
                    فایلی برای این نامه ثبت نشده است.
                  </div>
                ) : (
                  <div className="space-y-2">
                    {attachments.map((attachment) => (
                      <div
                        key={attachment.id}
                        className="flex flex-wrap items-center justify-between gap-3 rounded-md border p-3 text-sm"
                      >
                        <div className="min-w-0">
                          <div className="break-words font-medium">
                            {attachment.file_name}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            {attachment.content_type} -{" "}
                            {formatFileSize(attachment.size_bytes)}
                          </div>
                        </div>

                        <div className="flex gap-2">
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => void handleOpenAttachment(attachment)}
                          >
                            مشاهده
                          </Button>

                          <Button
                            type="button"
                            variant="destructive"
                            size="sm"
                            disabled={deletingAttachmentID === attachment.id}
                            onClick={() =>
                              void handleDeleteAttachment(attachment)
                            }
                          >
                            {deletingAttachmentID === attachment.id
                              ? t.commonLoading
                              : t.commonDelete}
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}

                <div className="mt-8 flex flex-wrap items-center justify-end gap-3 border-t pt-6">
                    <Button
                        type="submit"
                        form="edit-letter-form"
                        disabled={saving}
                    >
                        {saving ? t.commonLoading : "ذخیره تغییرات"}
                    </Button>

                    <Button
                        type="button"
                        variant="outline"
                        onClick={() => router.push("/letters")}
                        disabled={saving}
                    >
                        {t.commonCancel}
                    </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </AppShell>
    </ProtectedRoute>
  );
}

function formatFileSize(sizeBytes: number) {
  if (sizeBytes < 1024) {
    return `${sizeBytes} B`;
  }

  if (sizeBytes < 1024 * 1024) {
    return `${(sizeBytes / 1024).toFixed(1)} KB`;
  }

  return `${(sizeBytes / 1024 / 1024).toFixed(1)} MB`;
}
