import React, { useContext } from "react";
import { Select, message } from "antd";
import { SetChain } from "../../wailsjs/go/main/App";
import { AppContext } from "../appcontext";
const { Option } = Select;

export const ChainSelector: React.FC = () => {
  const { chainState, setChainState, address, status } = useContext(AppContext);

  const defaultChain = chainState.chain || "mainnet";
  const changeChain = (value: string) => {
    if (status === "Loading...") {
      message.warning(
        "Please wait for the current operation to finish or press ESC."
      );
      return;
    }
    setChainState({ ...chainState, chain: value });
    SetChain(value, address);
  };

  return (
    <Select defaultValue={defaultChain} onChange={changeChain}>
      <Option value="mainnet">mainnet</Option>
      <Option value="optimism">optimism</Option>
      <Option value="polygon">polygon</Option>
      <Option value="sepolia">sepolia</Option>
    </Select>
  );
};
