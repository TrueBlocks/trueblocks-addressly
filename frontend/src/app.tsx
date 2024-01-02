import React, { useState, useEffect } from "react";
import { Export, Reload } from "../wailsjs/go/main/App";
import * as Wails from "../wailsjs/runtime";
import logo from "./assets/images/logo.png";
import { BorderedWord } from "./components";
import "./App.css";
import { BarChart } from "./barchart";

export const App: React.FC = () => {
  const [address, setAddress] = useState("trueblocks.eth");
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
    <div className="panel window">
      <div className="panel header">
        <Logo />
        <Title />
        <ChainSelector />
      </div>
      <div className="panel main">
        <SideBar
          address={address}
          setAddress={setAddress}
          exportTxs={exportTxs}
          reloadTxs={reloadTxs}
        />
        <InnerPanel status={status} />
      </div>
      <div className="panel footer">
        <Config />
        <Copyright />
        <IconTray />
      </div>
    </div>
  );
};

var Logo = function () {
  return (
    <div className="panel header-left">
      <img className="logo" src={logo} alt="logo" />
    </div>
  );
};

var Title = function () {
  return <div className="panel header-middle">Account Browser</div>;
};

var ChainSelector = function () {
  var [price, setPrice] = useState(0.0);
  var [latest, setLatest] = useState(0);
  useEffect(() => {
    const update = (price: number) => {
      if (price > 0.0) {
        setPrice(price);
      }
    };
    Wails.EventsOn("price", update);
    return () => {
      Wails.EventsOff("price");
    };
  }, []);
  useEffect(() => {
    const update = (latest: number) => {
      if (latest > 0) {
        setLatest(latest);
      }
    };
    Wails.EventsOn("latest", update);
    return () => {
      Wails.EventsOff("latest");
    };
  }, []);

  return (
    <div className="panel header-right">
      <div className="price">{latest > 0 ? "latest: " + latest : ""}</div>
      <div className="price">{price > 0.0 ? "Eth price: " + price : ""}</div>
      <select id="chain-select">
        <option value="mainnet">Mainnet</option>
        <option value="optimism">Optimism</option>
        <option value="optimism">Sepolia</option>
        <option value="optimism">Polygon</option>
      </select>
    </div>
  );
};

interface SideBarProps {
  address: string;
  setAddress: React.Dispatch<React.SetStateAction<string>>;
  exportTxs: () => Promise<void>;
  reloadTxs: () => Promise<void>;
}

const SideBar: React.FC<SideBarProps> = ({
  address,
  setAddress,
  exportTxs,
  reloadTxs,
}) => {
  return (
    <div className="panel main-sidebar">
      <div>Address or ENS:</div>
      <div>
        <input
          className="input"
          onChange={(e) => setAddress(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && exportTxs()}
          value={address}
          placeholder="trueblocks.eth"
          autoComplete="off"
          name="input"
          autoFocus
        />
        <br />
        <button className="btn" onClick={reloadTxs} disabled={address === ""}>
          Export
        </button>
      </div>
    </div>
  );
};

interface StatusProps {
  status: string;
}

const InnerPanel: React.FC<StatusProps> = ({ status }) => {
  return (
    <div className="panel inner-panel">
      <Inner />
      <Logger status={status} />
    </div>
  );
};

const Inner = function () {
  var [info, setInfo] = useState("This is the info panel.");
  var [years, setYears] = useState("This is the years panel");
  var [months, setMonths] = useState("This is the months panel");
  var [toCount, setToCount] = useState("This is the toCount panel");
  var [fromCount, setFromCount] = useState("This is the fromCount panel");
  var [fromTopTen, setFromTopTen] = useState("This is the fromTopTen panel");
  var [toTopTen, setToTopTen] = useState("This is the toTopTen panel");

  useEffect(() => {
    const update = (info: string) => {
      setInfo(info);
    };
    Wails.EventsOn("info", update);
    return () => {
      Wails.EventsOff("info");
    };
  }, []);

  useEffect(() => {
    const update = (months: string) => {
      setMonths(months);
    };
    Wails.EventsOn("months", update);
    return () => {
      Wails.EventsOff("months");
    };
  }, []);

  useEffect(() => {
    const update = (years: string) => {
      setYears(years);
    };
    Wails.EventsOn("years", update);
    return () => {
      Wails.EventsOff("years");
    };
  }, []);

  useEffect(() => {
    const update = (toCount: string) => {
      setToCount(toCount);
    };
    Wails.EventsOn("toCount", update);
    return () => {
      Wails.EventsOff("toCount");
    };
  }, []);

  useEffect(() => {
    const update = (fromCount: string) => {
      setFromCount(fromCount);
    };
    Wails.EventsOn("fromCount", update);
    return () => {
      Wails.EventsOff("fromCount");
    };
  }, []);

  useEffect(() => {
    const update = (fromTopTen: string) => {
      setFromTopTen(fromTopTen);
    };
    Wails.EventsOn("fromTopTen", update);
    return () => {
      Wails.EventsOff("fromTopTen");
    };
  }, []);

  useEffect(() => {
    const update = (toTopTen: string) => {
      setToTopTen(toTopTen);
    };
    Wails.EventsOn("toTopTen", update);
    return () => {
      Wails.EventsOff("toTopTen");
    };
  }, []);

  let stringArray: string[] = [
    info,
    years,
    months,
    fromTopTen,
    fromCount,
    toTopTen,
    toCount,
  ];

  var info1 = function () {
    // var p = <CommaSeparatedDivs str={info} />;
    return (
      <div className="panel inner-panel-body-single">
        <BorderedWord word="info" content={info} />
      </div>
    );
  };

  var dates = function () {
    return (
      <div className="panel inner-panel-body-double">
        <BorderedWord
          word="Apps per Year"
          content={<BarChart dataString={years} />}
        />
        <BorderedWord
          word="Apps per Month"
          content={<BarChart dataString={months} />}
        />
      </div>
    );
  };

  var from = (
    <div className="panel inner-panel-body-double">
      <BorderedWord word="fromTopTen" content={fromTopTen} />
      <BorderedWord word="fromCount" content={fromCount} />
    </div>
  );

  var to = (
    <div className="panel inner-panel-body-double">
      <BorderedWord word="toTopTen" content={toTopTen} />
      <BorderedWord word="toCount" content={toCount} />
    </div>
  );

  var addrs = (
    <div className="panel inner-panel-body-double">
      {from}
      {to}
    </div>
  );

  return (
    <div className="inner-panel-body">
      {info1()}
      {dates()}
      {addrs}
    </div>
  );
};

type Props = {
  str: string;
};

const CommaSeparatedDivs: React.FC<Props> = ({ str }) => {
  const columns = str.split(",");
  return (
    <div>
      {columns.map((column, index) => (
        <div key={index}>{column.trim()}</div>
      ))}
    </div>
  );
};

const Logger: React.FC<StatusProps> = ({ status }) => {
  var [progress, setProgress] = useState("log messages are here");
  var [error, setError] = useState("");
  useEffect(() => {
    const update = (progress: string) => {
      setProgress(progress);
    };
    Wails.EventsOn("progress", update);
    return () => {
      Wails.EventsOff("progress");
    };
  }, []);

  useEffect(() => {
    const update = (error: string) => {
      setError(error);
    };
    Wails.EventsOn("error", update);
    return () => {
      Wails.EventsOff("error");
    };
  }, []);

  var classStr = "panel inner-panel-footer";
  var content = progress;
  if (error != "") {
    classStr += " error";
    content = error;
  } else if (status != "Loading...") {
    classStr += " empty";
    content = "Enter new addresses to the left...";
  }

  return (
    <div className={classStr}>
      <pre>{content}</pre>
    </div>
  );
};

var Config = function () {
  return <div className="panel footer-left">Config</div>;
};

var Copyright = function () {
  return (
    <div className="panel footer-middle">
      © 2024 TrueBlocks. All rights reserved.
    </div>
  );
};

var IconTray = function () {
  return (
    <div className="panel footer-right">
      <a
        href="https://twitter.com/trueblocks"
        target="_blank"
        rel="noopener noreferrer"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentcolor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M23 3a10.9 10.9.0 01-3.14 1.53 4.48 4.48.0 00-7.86 3v1A10.66 10.66.0 013 4s-4 9 5 13a11.64 11.64.0 01-7 2c9 5 20 0 20-11.5a4.5 4.5.0 00-.08-.83A7.72 7.72.0 0023 3z"></path>
        </svg>
      </a>
      <a
        href="http://yourwebsite.com"
        target="_blank"
        rel="noopener noreferrer"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentcolor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37.0 00-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44.0 0020 4.77 5.07 5.07.0 0019.91 1S18.73.65 16 2.48a13.38 13.38.0 00-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07.0 005 4.77 5.44 5.44.0 003.5 8.55c0 5.42 3.3 6.61 6.44 7A3.37 3.37.0 009 18.13V22"></path>
        </svg>
      </a>
      <a href="http://github.com" target="_blank" rel="noopener noreferrer">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          height="20"
          viewBox="-1 0 26 26"
          fill="none"
          stroke="currentcolor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M20.317 4.3698a19.7913 19.7913.0 00-4.8851-1.5152.0741.0741.0 00-.0785.0371c-.211.3753-.4447.8648-.6083 1.2495-1.8447-.2762-3.68-.2762-5.4868.0-.1636-.3933-.4058-.8742-.6177-1.2495a.077.077.0 00-.0785-.037 19.7363 19.7363.0 00-4.8852 1.515.0699.0699.0 00-.0321.0277C.5334 9.0458-.319 13.5799.0992 18.0578a.0824.0824.0 00.0312.0561c2.0528 1.5076 4.0413 2.4228 5.9929 3.0294a.0777.0777.0 00.0842-.0276c.4616-.6304.8731-1.2952 1.226-1.9942a.076.076.0 00-.0416-.1057c-.6528-.2476-1.2743-.5495-1.8722-.8923a.077.077.0 01-.0076-.1277c.1258-.0943.2517-.1923.3718-.2914a.0743.0743.0 01.0776-.0105c3.9278 1.7933 8.18 1.7933 12.0614.0a.0739.0739.0 01.0785.0095c.1202.099.246.1981.3728.2924a.077.077.0 01-.0066.1276 12.2986 12.2986.0 01-1.873.8914.0766.0766.0 00-.0407.1067c.3604.698.7719 1.3628 1.225 1.9932a.076.076.0 00.0842.0286c1.961-.6067 3.9495-1.5219 6.0023-3.0294a.077.077.0 00.0313-.0552c.5004-5.177-.8382-9.6739-3.5485-13.6604a.061.061.0 00-.0312-.0286zM8.02 15.3312c-1.1825.0-2.1569-1.0857-2.1569-2.419.0-1.3332.9555-2.4189 2.157-2.4189 1.2108.0 2.1757 1.0952 2.1568 2.419.0 1.3332-.9555 2.4189-2.1569 2.4189zm7.9748.0c-1.1825.0-2.1569-1.0857-2.1569-2.419.0-1.3332.9554-2.4189 2.1569-2.4189 1.2108.0 2.1757 1.0952 2.1568 2.419.0 1.3332-.946 2.4189-2.1568 2.4189z"></path>
        </svg>
      </a>
      <a
        href="https://trueblocks.io/"
        target="_blank"
        rel="noopener noreferrer"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          height="20"
          viewBox="0 0 420 420"
          stroke="currentcolor"
          fill="none"
        >
          <path strokeWidth="26" d="M209,15a195,195 0 1,0 2,0z"></path>
          <path
            strokeWidth="18"
            d="m210,15v390m195-195H15M59,90a260,260 0 0,0 302,0 m0,240 a260,260 0 0,0-302,0M195,20a250,250 0 0,0 0,382 m30,0 a250,250 0 0,0 0-382"
          ></path>
        </svg>
      </a>
    </div>
  );
};
