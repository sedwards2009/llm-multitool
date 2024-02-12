import { ChangeEvent, SyntheticEvent, useRef } from "react";
import { AttachedFile } from "./data";
import { FileAttachmentsList } from "./fileattachmentslist";
import { isImage } from "./mimetype_utils";

export interface Props {
  sessionId: string;
  attachedFiles: AttachedFile[];
  onUploadFile: (file: File) => void;
  onDeleteFile?: (filename: string) => void | null | undefined;
}

export function FileAttachments({sessionId, attachedFiles, onUploadFile, onDeleteFile}: Props): JSX.Element {

  const fileInputRef = useRef<HTMLInputElement>(null);
  const onClick = async () => {
    const current = fileInputRef.current;
    if (current != null) {
      current.click();
    }
  };

  const onFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files == null || files.length === 0) {
      return;
    }
    onUploadFile(files[0]);
  };

  const onPaste = async (event: SyntheticEvent) => {
    event.preventDefault();
    const e = event.nativeEvent as ClipboardEvent;
    const files = e.clipboardData?.files;
    if (files == null) {
      return;
    }

    for (const clipboardItem of files) {
      if (isImage(clipboardItem.type)) {
        onUploadFile(clipboardItem);
      }
    }
  };

  return <div>
    <div className="gui-packed-row">
      <button className="small compact" onClick={onClick}>
        <i className="fas fa-paperclip"></i>
        {" Attach File"}
      </button>

      <input
        type="text"
        className="char-max-width-12"
        placeholder="Ctrl+V Paste Files Here"
        onPaste={onPaste}
        readOnly={true}
      />

      <form>
        <input type="file" ref={fileInputRef} className="hidden" onChange={onFileChange} />
      </form>
    </div>

    <FileAttachmentsList
      sessionId={sessionId}
      attachedFiles={attachedFiles}
      onDeleteFile={onDeleteFile}
    />
  </div>;
}

