import React from "react";

type LabelProps = Omit<Rect.ComponentPropsWithoutRef<"label">, "className">;

const Label = (props: LabelProps) => {
  return <label {...props} className="block font-bold text-lg" />;
};

export default Label;
