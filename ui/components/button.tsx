import React from "react";

type ButtonProps = Omit<React.ComponentPropsWithRef<"button">, "className">;

const Button = React.forwardRef((props: ButtonProps, ref) => {
  return (
    <button
      {...props}
      ref={ref}
      className="px-4 py-2 border border-pink-600 hover:bg-pink-600 active:bg-pink-600 rounded font-bold"
    />
  );
});

export default Button;
