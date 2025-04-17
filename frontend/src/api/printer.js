export const printer = {
  list,
};

async function list() {
  return [
    {
      id: 1,
      name: "HP",
      status: "online",
    },
    {
      id: 2,
      name: "Epson",
      status: "offline",
    },
  ];
}
