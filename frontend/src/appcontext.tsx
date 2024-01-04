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
  chainState: IChainState;
  setChainState: React.Dispatch<React.SetStateAction<IChainState>>;
}

export const AppContext = createContext<IAppContext>({
  address: "",
  setAddress: () => {},
  info: "",
  setInfo: () => {},
  current: 0,
  setCurrent: () => {},
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

  const value: IAppContext = {
    address,
    setAddress,
    info,
    setInfo,
    current,
    setCurrent,
    chainState,
    setChainState,
  };

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}
