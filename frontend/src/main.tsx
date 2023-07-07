import ReactDOM from "react-dom/client";
import { App } from "./app.tsx";
import "./theme/main.scss";
import "./theme/fonts/fontawesome-fontface.scss";
import "./theme/fonts/font-awesome.scss";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <App />
);
