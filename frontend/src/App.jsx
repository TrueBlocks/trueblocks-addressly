import React, { useState } from "react";
import logo from "./assets/images/logo.png";
import { Export } from "../wailsjs/go/main/App";
import styled from "styled-components";
import * as C from './components'

function App() {
  const [name, setName] = useState("");
  const [status, setStatus] = useState("");
  const [loading, setLoading] = useState(false);

  const exportTxs = async () => {
    setLoading(true);
    setStatus("");
    const result = await Export(name, "");
    setStatus(result);
    setLoading(false);
  };

  return (
    <C.StyledApp>
      <C.Header>
        <C.HeaderSide width="25%"></C.HeaderSide>
        <C.HeaderMiddle width="50%">
          TrueBlocks
          <br />
          Account Explorer
          <C.Logo src={logo} alt="logo" />
        </C.HeaderMiddle>
        <C.HeaderSide width="25%">
          <select id="chain-select">
            <option value="mainnet">Mainnet</option>
            <option value="optimism">Optimism</option>
          </select>
        </C.HeaderSide>
      </C.Header>
      <div>
        <C.Prompt>Enter an address or ENS name below 👇</C.Prompt>
        <C.InputBox>
          <input
            className="input"
            onChange={(e) => setName(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && exportTxs()}
            value={name}
            placeholder="trueblocks.eth"
            autoComplete="off"
            name="input"
            autoFocus
          />
          <button
            className="btn"
            onClick={exportTxs}
            disabled={loading || name === ""}
          >
            Export
          </button>
        </C.InputBox>
        <C.Result>{loading ? "Loading..." : status}</C.Result>
      </div>
    </C.StyledApp>
  );
}

export default App;

