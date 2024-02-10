import { AttachedFile } from "./data";
import { fileURL } from "./dataloading";

export interface Props {
  sessionId: string;
  attachedFiles: AttachedFile[];
  onDeleteFile?: (filename: string) => void | null | undefined;
}

function isImage(mimeType: string): boolean {
  return ['image/png', 'image/jpeg', 'image/gif'].includes(mimeType);
}

export function FileAttachmentsList({sessionId, attachedFiles, onDeleteFile}: Props): JSX.Element {
  return <div>
    {
      attachedFiles.map(af => {
        const filename = af.filename;
        const deleteButton = onDeleteFile &&
          <button className="microtool danger" onClick={() => onDeleteFile(filename)}><i className="fa fa-times"></i></button>;
        if (isImage(af.mimeType)) {
          return <div key={filename} className="uploaded-image">
            <img src={fileURL(sessionId, af.filename)} />
            {deleteButton}
          </div>;
        } else {
          return <div key={filename}>
            {af.originalFilename}
            {deleteButton}
          </div>;
        }
      })
    }
  </div>;
}
