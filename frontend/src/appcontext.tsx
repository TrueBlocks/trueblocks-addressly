import React, { useEffect, createContext, useState, ReactNode } from "react";
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
  monitors: string;
  setMonitors: React.Dispatch<React.SetStateAction<string>>;
  chainState: IChainState;
  setChainState: React.Dispatch<React.SetStateAction<IChainState>>;
}

var defaultMonitors = `trueblocks.eth
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
  monitors: defaultMonitors,
  setMonitors: () => {},
  chainState: { block: "", date: "", price: "", chain: "" },
  setChainState: () => {},
});

interface AppProviderProps {
  children: ReactNode;
}

export function AppProvider({ children }: AppProviderProps) {
  const [address, setAddress] = useState("trueblocks.eth");
  const [current, setCurrent] = useState(0);
  const [info, setInfo] = useState("");
  const [monitors, setMonitors] = useState(defaultMonitors);
  const [chainState, setChainState] = useState<IChainState>({
    block: "",
    date: "",
    price: "",
    chain: "",
  });

  useEffect(() => {
    const update = (newBlock: string) => {
      var parts = newBlock.split("|");
      setChainState({
        block: parts[0],
        date: parts[1],
        price: parts[2],
        chain: parts[3],
      });
    };
    Wails.EventsOn("chainState", update);
    return () => {
      Wails.EventsOff("chainState");
    };
  }, []);

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
    const update = (monitors: string) => {
      setMonitors(monitors);
    };
    Wails.EventsOn("monitors", update);
    return () => {
      Wails.EventsOff("monitors");
    };
  }, []);

  const value: IAppContext = {
    address,
    setAddress,
    info,
    setInfo,
    current,
    setCurrent,
    monitors,
    setMonitors,
    chainState,
    setChainState,
  };

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}
