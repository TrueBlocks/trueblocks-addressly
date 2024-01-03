import React, { useState, useEffect } from "react";
import * as Wails from "../../wailsjs/runtime";
import "../App.css";

export const Logger = () => {
  const [progress, setProgress] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    const updateProgress = (progress: string) => setProgress(progress);
    Wails.EventsOn("progress", updateProgress);

    const updateError = (error: string) => setError(error);
    Wails.EventsOn("error", updateError);

    return () => {
      Wails.EventsOff("progress");
      Wails.EventsOff("error");
    };
  }, []);

  const classStr = `panel inner-panel-footer ${error ? "error" : progress ? "" : "empty"}`;
  const content = error || progress || "Enter new addresses to the left...";

  return (
    <div className={classStr}>
      <pre>{content}</pre>
    </div>
  );
};
