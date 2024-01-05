import React, { useState, useContext, useEffect } from "react";
import * as Wails from "../../wailsjs/runtime";
import { Input, Button, Typography, Space, Switch } from "antd";
import { AppContext } from "../appcontext";
import { Export, Freshen, SetExportExcel } from "../../wailsjs/go/main/App";
import { ChainSelector } from "./chainselector";
const { Text } = Typography;

export const SideBar: React.FC = () => {
  const { address, setAddress } = useContext(AppContext);
  const [status, setStatus] = useState("Enter an address at the left...");
  const [exportExcel, setExportExcel] = useState(false);

  const mode = "";
  const exportTxs = async () => {
    if (status == "Loading...") {
      return;
    }
    setStatus("Loading...");
    await Export(address, mode);
    setStatus("");
  };

  const reloadTxs = async () => {
    if (status == "Loading...") {
      return;
    }
    setStatus("Loading...");
    await Freshen(address, mode);
    setStatus("");
  };

  useEffect(() => {
    const update = (exportExcel: string) => {
      setExportExcel(exportExcel === "true" ? true : false);
    };
    Wails.EventsOn("exportExcel", update);
    return () => {
      Wails.EventsOff("exportExcel");
    };
  }, []);

  const handleToggle = (value: boolean) => {
    setExportExcel(value);
    SetExportExcel(value);
  };

  return (
    <Space
      direction="vertical"
      size="middle"
      style={{ textAlign: "left", marginTop: 20 }}
    >
      <Text style={{ textAlign: "left", color: "white", fontSize: ".9em" }}>
        Address or ENS:
      </Text>
      <Input
        onChange={(e) => setAddress(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && exportTxs()}
        value={address}
        placeholder="trueblocks.eth"
        autoFocus
        style={{ textAlign: "left" }}
      />
      <Text style={{ textAlign: "left", color: "white", fontSize: ".9em" }}>
        Chain:
      </Text>
      <ChainSelector />
      <Text style={{ color: "white" }}>
        Export to Excel:{" "}
        <Switch
          checked={exportExcel}
          onClick={handleToggle}
          size="small"
        ></Switch>
      </Text>
      <Button
        onClick={reloadTxs}
        disabled={address === ""}
        type="primary"
        style={{ textAlign: "left" }}
      >
        {exportExcel ? "Export to Excel" : "Freshen"}
      </Button>
    </Space>
  );
};
