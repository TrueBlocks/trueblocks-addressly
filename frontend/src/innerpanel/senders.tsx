import React, { useState, useEffect, useContext } from "react";
import { AppContext } from "../appcontext";
import * as Wails from "../../wailsjs/runtime";
import { TransactionsTable } from "./tx_line";
import "./inner.css";

export const Senders: React.FC = () => {
  const [asSender, setSender] = useState("This is the asSender panel");
  const { info } = useContext(AppContext);
  const currAddr = info.split(",")[1];

  useEffect(() => {
    const update = (asSender: string): void => {
      setSender(asSender);
    };
    Wails.EventsOn("asSender", update);
    return () => {
      Wails.EventsOff("asSender");
    };
  }, []);

  const columns = asSender.replace("/,$/g", "").split(",");
  return <TransactionsTable data={columns} currAddr={currAddr} />
};
