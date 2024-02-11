import { } from "react";
import { isImage } from "./mimetype_utils";

export interface Props {
  fileUrl: string;
  mimetype: string | null;
  onCloseRequest: () => void;
}

export function FilePreview({fileUrl, mimetype, onCloseRequest}: Props): JSX.Element {
  return (
    <div className="file-preview" onClick={onCloseRequest}>
      {isImage(mimetype) && <img src={fileUrl} />}
      {!isImage(mimetype) && <p>${fileUrl}</p>}
    </div>
  );
}
