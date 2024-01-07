import React, { useState, useContext, useEffect } from "react";
import * as Wails from "../../wailsjs/runtime";
import { Input, Button, Typography, Space, Switch, message } from "antd";
import { AppContext } from "../appcontext";
import { Export, Freshen, SetExportExcel } from "../../wailsjs/go/main/App";
import { ChainSelector } from "./chainselector";
const { Text } = Typography;

export const SideBar: React.FC = () => {
  const { address, setAddress, status, setStatus } = useContext(AppContext);
  const [exportExcel, setExportExcel] = useState(false);

  const mode = "";
  const exportTxs = async (): Promise<void> => {
    if (status == "Loading...") {
      message.warning(
        "Please wait for the current operation to finish or press ESC."
      );
      return;
    }
    setStatus("Loading...");
    await Export(address, mode);
    setStatus("");
  };

  const reloadTxs = async (): Promise<void> => {
    if (status == "Loading...") {
      message.warning(
        "Please wait for the current operation to finish or press ESC."
      );
      return;
    }
    setStatus("Loading...");
    await Freshen(address, mode);
    setStatus("");
  };

  useEffect(() => {
    const update = (exportExcel: string): void => {
      setExportExcel(exportExcel === "true");
    };
    Wails.EventsOn("exportExcel", update);
    return () => {
      Wails.EventsOff("exportExcel");
    };
  }, []);

  const handleToggle = (value: boolean): void => {
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
        onChange={(e) => {
          setAddress(e.target.value);
        }}
        onKeyDown={async (e) => await (e.key === "Enter" && exportTxs())}
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
