import { navigate, usePath } from "raviger";

export interface Props {}

export function TitleBar({ }: Props): JSX.Element {
  const path = usePath();

  const onSettingsClicked = () => {
    console.log(`path: ${path}`);
    navigate("/settings", {query: {from: path}});
  };

  return <div className="gui-packed-row width-100pc">
    <h1 className="expand">LLM Workbench</h1>
    <button className="small compact" onClick={onSettingsClicked}>
      <i className="fas fa-cog"></i>
      {" Settings"}
    </button>
  </div>;
}
