import React, { useEffect, createContext, useState, ReactNode } from "react";
import * as Wails from "../wailsjs/runtime";

interface IChainState {
  block: string;
  date: string;
  price: string;
}

interface IAppContext {
  address: string;
  setAddress: React.Dispatch<React.SetStateAction<string>>;
  chainState: IChainState;
  setChainState: React.Dispatch<React.SetStateAction<IChainState>>;
}

export const AppContext = createContext<IAppContext>({
  address: "",
  setAddress: () => {},
  chainState: { block: "", date: "", price: "" },
  setChainState: () => {},
});

interface AppProviderProps {
  children: ReactNode;
}

export function AppProvider({ children }: AppProviderProps) {
  const [address, setAddress] = useState("trueblocks.eth");
  const [chainState, setChainState] = useState<IChainState>({
    block: "",
    date: "",
    price: "",
  });

  useEffect(() => {
    const update = (newBlock: string) => {
      var parts = newBlock.split("|");
      setChainState({
        block: parts[0],
        date: parts[1],
        price: parts[2],
      });
    };
    Wails.EventsOn("chainState", update);
    return () => {
      Wails.EventsOff("chainState");
    };
  }, []);

  const value: IAppContext = {
    address,
    setAddress,
    chainState,
    setChainState,
  };

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}
