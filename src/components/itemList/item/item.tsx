import React, { useState } from "react";

import "./item.css";

interface Props {
  title: string;
  id: string;
  onRemove: (id: string) => void;
}

export default function Item(props: Props) {
  const [hover, setHover] = useState(false);
  return (
    <div
      className="item-container"
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      onClick={() => props.onRemove(props.id)}
    >
      {hover ? "Delete" : props.title}
    </div>
  );
}
