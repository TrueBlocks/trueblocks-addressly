import React, { useState, useContext } from "react";
import { ChangeName } from "../../wailsjs/go/main/App";
import { Typography, Input, Form } from "antd";
import { AppContext } from "../appcontext";

type EditableAddressProps = {
  name: string;
};

export const EditableAddress: React.FC<EditableAddressProps> = ({ name }) => {
  const [editable, setEditable] = useState<boolean>(false);
  const [textValue, setTextValue] = useState<string>(name);

  var parts = name.split("|");
  name = parts[0];
  var address = parts[1];

  const handleEdit = (): void => {
    // setTextValue(name);
    // setEditable(true);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    setTextValue(e.target.value);
  };

  const handleEnter = (e: React.KeyboardEvent<HTMLInputElement>): void => {
    if (e.key === "Enter") {
      ChangeName(textValue, address);
      handleSubmit();
    }
  };

  const handleSubmit = (): void => {
    setEditable(false);
  };

  return (
    <div>
      {editable ? (
        <Form.Item>
          <Input
            value={textValue}
            onChange={handleChange}
            onPressEnter={handleEnter}
            onBlur={handleSubmit}
            autoFocus
          />
        </Form.Item>
      ) : (
        <Typography.Text
          // style={{ backgroundColor: "pink" }} //useContext(AppContext).theme.background}}
          onClick={handleEdit}
        >
          {name}
        </Typography.Text>
      )}
    </div>
  );
};
