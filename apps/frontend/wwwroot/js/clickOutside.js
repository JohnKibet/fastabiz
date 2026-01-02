window.miniCartClickOutside = {
  _map: {},

  // register(dotNetRef) => returns token string
  register: function (dotNetHelper) {
    const id = Math.random().toString(36).slice(2);

    const handler = function (event) {
      const dropdown = document.getElementById('mini-cart-dropdown');
      if (!dropdown) return; // nothing to do if dropdown isn't rendered
      if (!dropdown.contains(event.target)) {
        // call the .NET method named 'CloseMiniCart'
        dotNetHelper.invokeMethodAsync('CloseMiniCart').catch(err => {
          console.error('miniCartClickOutside invoke failed', err);
        });
      }
    };

    document.addEventListener('mousedown', handler);
    this._map[id] = handler;
    return id;
  },

  unregister: function (token) {
    const handler = this._map[token];
    if (handler) {
      document.removeEventListener('mousedown', handler);
      delete this._map[token];
    }
  }
};
