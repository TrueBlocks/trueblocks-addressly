import React, { useState, useEffect, useContext } from "react";
import * as Wails from "../../wailsjs/runtime";
import "./inner.css";
import { TransactionsTable } from "./tx_line";
import { AppContext } from "../appcontext";

export const Recipients: React.FC = () => {
  const [asRecipient, setRecipient] = useState("This is the asRecipient panel");
  const { info } = useContext(AppContext);
  const currAddr = info.split(",")[1];

  useEffect(() => {
    const update = (asRecipient: string): void => {
      setRecipient(asRecipient);
    };
    Wails.EventsOn("asRecipient", update);
    return () => {
      Wails.EventsOff("asRecipient");
    };
  }, []);

  const columns = asRecipient.split(",");
  return <TransactionsTable data={columns} currAddr={currAddr} />;
};
