import { ChangeEvent, useRef } from "react";
import { AttachedFile } from "./data";
import { FileAttachmentsList } from "./fileattachmentslist";

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

  return <div>
    <div className="gui-packed-row">
      <button className="small compact" onClick={onClick}>
        <i className="fas fa-paperclip"></i>
        {" Attach File"}
      </button>
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

