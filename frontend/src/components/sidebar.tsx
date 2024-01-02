import React from "react";
import { Input, Button, Typography, Space } from "antd";
const { Text } = Typography;

interface SideBarProps {
  address: string;
  setAddress: React.Dispatch<React.SetStateAction<string>>;
  exportTxs: () => Promise<void>;
  reloadTxs: () => Promise<void>;
}

export const SideBar: React.FC<SideBarProps> = ({
  address,
  setAddress,
  exportTxs,
  reloadTxs,
}) => {
  return (
    <Space direction="vertical" size="middle" style={{ marginTop: 20 }}>
      <Text>Address or ENS:</Text>
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
