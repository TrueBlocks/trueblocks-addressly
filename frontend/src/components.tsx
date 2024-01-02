import React from "react";

interface BorderedWordProps {
  word: string;
  content: any;
}

export const BorderedWord: React.FC<BorderedWordProps> = ({
  word,
  content,
}) => {
  return (
    <div className="bordered-word-container">
      <div className="bordered-word">{word}</div>
      <div>{content}</div>
    </div>
  );
};
