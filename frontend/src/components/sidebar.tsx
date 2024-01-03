import React, { useContext } from "react";
import { Input, Button, Typography, Space } from "antd";
const { Text } = Typography;
import { AppContext } from "../appcontext";

interface SideBarProps {
  exportTxs: () => Promise<void>;
  reloadTxs: () => Promise<void>;
}

export const SideBar: React.FC<SideBarProps> = ({ exportTxs, reloadTxs }) => {
  var { address, setAddress } = useContext(AppContext);
  return (
    <Space direction="vertical" size="middle" style={{ marginTop: 20 }}>
      <Text style={{ textAlign: "left", color: "white", fontSize: ".9em" }}>
        Address or ENS:
      </Text>
      <Input
        onChange={(e) => setAddress(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && exportTxs()}
        value={address}
        placeholder="trueblocks.eth"
        autoFocus
      />
      <Button onClick={reloadTxs} disabled={address === ""} type="primary">
        Reload
      </Button>
    </Space>
  );
};
