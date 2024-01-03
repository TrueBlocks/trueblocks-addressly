import React, { useState, useContext } from "react";
import { Input, Button, Typography, Space, Switch } from "antd";
import { AppContext } from "../appcontext";
import { Export, Reload } from "../../wailsjs/go/main/App";
import { ChainSelector } from "./chainselector";
const { Text } = Typography;

export const SideBar: React.FC = () => {
  const { address, setAddress } = useContext(AppContext);
  const [status, setStatus] = useState("Enter an address at the left...");

  const mode = "";
  const exportTxs = async () => {
    if (status == "Loading...") {
      return;
    }
    setStatus("Loading...");
    await Export(address, mode, false);
    setStatus("");
  };

  const reloadTxs = async () => {
    if (status == "Loading...") {
      return;
    }
    setStatus("Loading...");
    await Reload(address, mode, false);
    setStatus("");
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
      <Button
        onClick={reloadTxs}
        disabled={address === ""}
        type="primary"
        style={{ textAlign: "left" }}
      >
        Reload
      </Button>
      <Text style={{ textAlign: "left", color: "white", fontSize: ".9em" }}>
        Chain:
      </Text>
      <ChainSelector />
      <Text style={{ color: "white" }}>
        Export to Excel: <Switch size="small"></Switch>
      </Text>
    </Space>
  );
};
