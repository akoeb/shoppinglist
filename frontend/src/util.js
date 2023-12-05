export const httpOptions = (m, obj) => {
  // this object will be returned
  var options = {
    method: m !== undefined ? m : "GET",
    mode: "cors",
    cache: "no-cache",
    credentials: "same-origin",
    headers: {
      "Content-Type": "application/json",
    },
    redirect: "follow",
  };
  if (obj !== undefined) {
    options["body"] = JSON.stringify(obj);
  }
  return options;
};

export const backend = (path) => {
  // pass in environment from ENV var
  let env = process.env.NODE_ENV;
  let url = "http://localhost:8080/";
  if (env === "production") {
    url = "https://koeb.me/list/";
  }
  console.log("call to backend: " + url + path + " (" + env + ")");
  return url + path;
};

export const uid = () => {
  return (
    Date.now().toString(36) +
    Math.random().toString(36).substring(2, 12).padStart(12, 0)
  );
};

export const reorderStore = (store, start, target) => {
  const newStore = store;

  //
  if (start < target) {
    newStore.items.splice(target + 1, 0, newStore.items[start]);
    newStore.items.splice(start, 1);
  } else {
    newStore.items.splice(target, 0, newStore.items[start]);
    newStore.items.splice(start + 1, 1);
  }

  // sync order numbers in elements accordingly
  newStore.items.forEach(function (x, idx) {
    x.orderno = idx;
  });
  newStore.local = true;
  return newStore;
};

// if the store changes through user actions, we sync local storage and backend
// if the store changes through remote actions, we only sync store without local storage
export function notifyWhenChanges(store, local, callback) {
  return store.subscribe((value) => {
    if (local) {
      callback(value);
    }
  });
}
