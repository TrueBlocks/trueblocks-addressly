import React, { useContext } from "react";
import { Select } from "antd";
import { SetChain } from "../../wailsjs/go/main/App";
import { AppContext } from "../appcontext";
const { Option } = Select;

export const ChainSelector: React.FC = () => {
  const { chainState, setChainState, address } = useContext(AppContext);

  const defaultChain = chainState.chain || "mainnet";
  const changeChain = (value: string) => {
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
