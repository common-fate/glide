import { render } from "react-dom";
import { Routes } from "./utils/generouted";

const container = document.getElementById("app")!;
render(<Routes />, container);
