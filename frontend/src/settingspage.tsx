import { navigate, useQueryParams } from "raviger";
import { ModelOverview } from "./data";
import { useState } from "react";


export interface Props {
  modelOverview: ModelOverview;
  rescanModels: () => Promise<void>;
}

export function SettingsPage({ modelOverview, rescanModels }: Props): JSX.Element {
  const [{from}, _] = useQueryParams();
  const onBackClicked = () => {
    navigate(from);
  };

  const [isScanning, setIsScanning] = useState<boolean>(false);
  const onRescanClicked = () => {
    (async () => {
      setIsScanning(true);
      await rescanModels();
      setIsScanning(false);
    })();
  };

  return (
    <>
      <h2>Settings</h2>

      <h3>Models</h3>
      <table>
        <thead>
          <th>Model</th>
          <th>Supports Continue</th>
          <th>Supports Replies</th>
          <th>Supports Images</th>
        </thead>
        <tbody>
        {
          modelOverview.models.map(m => {
            return <tr key={m.id}>
              <td>{m.name}</td>
              <td>{m.supportsContinue ? "Yes" : "No"}</td>
              <td>{m.supportsReply ? "Yes" : "No"}</td>
              <td>{m.supportsImages ? "Yes" : "No"}</td>
            </tr>;
          })
        }
        </tbody>
      </table>

      <button className="success" onClick={onRescanClicked} disabled={isScanning}>
        { isScanning && <i className="fas fa-spinner fa-spin"></i>}
        {" Scan Models"}
      </button>
      <br/ >
      <br/ >
      <button className="small primary" onClick={onBackClicked}>
        <i className="fas fa-arrow-left"></i>
        {" Back"}
      </button>
    </>
  );
}
