import { SyntheticEvent, useContext } from "react";
import { AttachedFile } from "./data";
import { fileURL } from "./dataloading";
import { PreviewRequestContext } from "./filepreviewrequestcontext";
import { isImage } from "./mimetype_utils";

export interface Props {
  sessionId: string;
  attachedFiles: AttachedFile[];
  onDeleteFile?: (filename: string) => void | null | undefined;
}


export function FileAttachmentsList({sessionId, attachedFiles, onDeleteFile}: Props): JSX.Element {
  const onPreviewFileUrlRequest = useContext(PreviewRequestContext);

  return <div>
    {
      attachedFiles.map(af => {
        const filename = af.filename;
        const deleteButton = onDeleteFile &&
          <button
            className="microtool danger"
            onClick={(event: SyntheticEvent) => {
              event.preventDefault();
              event.stopPropagation();
              onDeleteFile(filename);
            }}
          >
            <i className="fa fa-times"></i>
          </button>;

        if (isImage(af.mimeType)) {
          return <div
            key={filename}
            className="uploaded-image"
            onClick={() => {
              if (onPreviewFileUrlRequest != null) {
                onPreviewFileUrlRequest(fileURL(sessionId, af.filename), af.mimeType);
              }
            }}
            >
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
