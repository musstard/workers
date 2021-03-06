require("./polyfill_performance.js");
require("./wasm_exec.js");

addEventListener("fetch", (event) => {
  event.respondWith(handleRequest(event.request));
});

async function handleRequest(request) {
  // Create our instance, with an imported function
  const go = new Go();
  go.importObject.env["command-line-arguments.sayHello"] = () => {
    console.log("Hello from the imported function!");
  };
  const instance = await WebAssembly.instantiate(WASM, go.importObject);

  go.run(instance);

  // Our Golang has access to the imported function
  const test = goGetter;
  //const test = MyGoFunc()
  console.log(test);
  
  return new Response(`0 x 0`, { status: 200 });
}
