import React, { useEffect } from "react";
import Button from "@material-ui/core/Button";
import TextField from "@material-ui/core/TextField";
import AddIcon from "@material-ui/icons/PostAdd";
import DeleteIcon from "@material-ui/icons/DeleteForever";

import "./controls.css";

interface Props {
  item: string;
  allItems: string[];
  onChange: (newItem: string) => void;
  onAdd: () => void;
  onClear: () => void;
}

export default function Controls(props: Props) {
  useEffect(() => {
    console.log("Controls updated");
    console.log("Item Value = ", props.item);
  }, [props.item]);

  return (
    <div className="controls-container">
      <TextField
        id="standard-basic"
        placeholder="Add a New Item"
        value={props.item}
        onChange={(e) => props.onChange(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            props.onAdd();
          }
        }}
      />
      <Button
        className="controls-button"
        variant="contained"
        color="primary"
        onClick={() => {
          props.onAdd();
        }}
        startIcon={<AddIcon />}
        disabled={!props.item}
      >
        Add
      </Button>
      <Button
        className="controls-button"
        variant="contained"
        color="secondary"
        onClick={() => {
          props.onClear();
        }}
        startIcon={<DeleteIcon />}
        disabled={props.allItems.length === 0}
      >
        Clear All
      </Button>
    </div>
  );
}
