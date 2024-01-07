import React, {
  useEffect,
  createContext,
  useState,
  type ReactNode
} from "react";
import * as Wails from "../wailsjs/runtime";

interface IChainState {
  chain: string;
  block: string;
  date: string;
  price: string;
}

interface IAppContext {
  address: string;
  setAddress: React.Dispatch<React.SetStateAction<string>>;
  info: string;
  setInfo: React.Dispatch<React.SetStateAction<string>>;
  current: number;
  setCurrent: React.Dispatch<React.SetStateAction<number>>;
  chartType: string;
  setChartType: React.Dispatch<React.SetStateAction<string>>;
  status: string;
  setStatus: React.Dispatch<React.SetStateAction<string>>;
  monitors: string;
  setMonitors: React.Dispatch<React.SetStateAction<string>>;
  chainState: IChainState;
  setChainState: React.Dispatch<React.SetStateAction<IChainState>>;
}

const defaultMonitors = `trueblocks.eth
vald.eth
meriam.eth
griff.eth
vitalik.eth
giveth.eth
ethereumfoundation.eth
gnosis.eth
makerdao.eth
molochdao.eth
0xdd94de9cfe063577051a5eb7465d08317d8808b6
ethdenver.eth
chasewright.eth
makingprogress.eth
omnianalytics.eth
austingriffith.eth`;

export const AppContext = createContext<IAppContext>({
  address: "",
  setAddress: () => {},
  info: "",
  setInfo: () => {},
  current: 0,
  setCurrent: () => {},
  chartType: "month",
  setChartType: () => {},
  status: "Enter an address at the left...",
  setStatus: () => {},
  monitors: defaultMonitors,
  setMonitors: () => {},
  chainState: { block: "", date: "", price: "", chain: "" },
  setChainState: () => {}
});

interface AppProviderProps {
  children: ReactNode;
}

export function AppProvider({ children }: AppProviderProps) {
  const [address, setAddress] = useState("trueblocks.eth");
  const [current, setCurrent] = useState(0);
  const [info, setInfo] = useState("");
  const [monitors, setMonitors] = useState(defaultMonitors);
  const [chartType, setChartType] = useState("month");
  const [status, setStatus] = useState("Enter an address at the left...");
  const [chainState, setChainState] = useState<IChainState>({
    block: "",
    date: "",
    price: "",
    chain: ""
  });

  useEffect(() => {
    const update = (cS: string): void => {
      const parts = cS.split("|");
      setChainState({
        block: parts[0],
        date: parts[1],
        price: parts[2],
        chain: parts[3]
      });
      // console.log("useEffect update chainState: ", cS, chainState);
    };
    Wails.EventsOn("chainState", update);
    return () => {
      Wails.EventsOff("chainState");
    };
  }, []);

  useEffect(() => {
    const update = (i: string): void => {
      setInfo(i);
      // console.log("useEffect update info: ", i, info);
    };
    Wails.EventsOn("info", update);
    return () => {
      Wails.EventsOff("info");
    };
  }, []);

  useEffect(() => {
    const update = (mon: string): void => {
      setMonitors(mon);
      // console.log("useEffect update monitors: ", mon, monitors);
    };
    Wails.EventsOn("monitors", update);
    return () => {
      Wails.EventsOff("monitors");
    };
  }, []);

  useEffect(() => {
    const update = (ct: string): void => {
      setChartType(ct);
      // console.log("useEffect update chartType: ", ct, chartType);
    };
    Wails.EventsOn("chartType", update);
    return () => {
      Wails.EventsOff("chartType");
    };
  }, []);

  const value: IAppContext = {
    address,
    setAddress,
    info,
    setInfo,
    current,
    setCurrent,
    chartType,
    setChartType,
    status,
    setStatus,
    monitors,
    setMonitors,
    chainState,
    setChainState
  };

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}
