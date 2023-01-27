import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import { Routes } from "./utils/generouted";

const root = ReactDOM.createRoot(document.getElementById("app") as HTMLElement);

root.render(
  <StrictMode>
    <Routes useErrorBoundary />
  </StrictMode>
);
