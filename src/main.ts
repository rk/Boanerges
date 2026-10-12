import { mount } from "svelte";
import App from "./App.svelte";
import "./style.css";

async function bootstrap(): Promise<void> {
  const params = new URLSearchParams(window.location.search);
  const e2e = params.get("e2eWelcome") === "1" ? "welcome" : import.meta.env.VITE_E2E === "1" ? "1" : "";
  if (e2e === "1" || e2e === "welcome") {
    const { installMygoMock } = await import("./e2e/mock-mygo");
    installMygoMock(e2e);
  }

  mount(App, { target: document.getElementById("app")! });
}

void bootstrap();
