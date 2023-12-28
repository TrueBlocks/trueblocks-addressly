import React, { useState } from "react";
import logo from "./assets/images/logo.png";
import { Export } from "../wailsjs/go/main/App";
import styled from "styled-components";
import {StyledApp, Logo,Prompt,InputBox,Result,Header,HeaderMiddle,HeaderSide,Body,BodyMiddle,BodySide,Footer,FooterMiddle,FooterSide}from './components'

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
    <StyledApp>
      <Header>
        <HeaderSide></HeaderSide>
        <HeaderMiddle>
          TrueBlocks
          <br />
          Account Explorer
          <Logo src={logo} alt="logo" />
        </HeaderMiddle>
        <HeaderSide>
          <select id="chain-select">
            <option value="mainnet">Mainnet</option>
            <option value="optimism">Optimism</option>
          </select>
        </HeaderSide>
      </Header>
      <Body>
        <BodySide>
          <Prompt>Enter an address or ENS</Prompt>
          <InputBox>
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
          </InputBox>
        </BodySide>
        <BodyMiddle>
          <Result>{loading ? "Loading..." : status}</Result>
        </BodyMiddle>
      </Body>
      <Footer>
        <FooterSide></FooterSide>
        <FooterMiddle>
          © 2024 TrueBlocks. All rights reserved.
        </FooterMiddle>
        <FooterSide></FooterSide>
      </Footer>
    </StyledApp>
  );
}

export default App;
