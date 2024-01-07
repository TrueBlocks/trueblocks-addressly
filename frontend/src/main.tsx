import React from "react";
import { createRoot } from "react-dom/client";
import "./style.css";
import { App } from "./app";
import { AppProvider } from "./appcontext";

const container = document.getElementById("root");
if (container != null) {
  const root = createRoot(container);
  root.render(
    <React.StrictMode>
      <AppProvider>
        <App />
      </AppProvider>
    </React.StrictMode>
  );
} else {
  console.error("Root element not found");
}
