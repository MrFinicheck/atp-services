import { api, clearToken, getToken, setToken, wailsError } from '../api/client';
import type { Role, User } from '../types';

let currentUser: User | null = null;
let currentView = 'dashboard';
let sidebarOpen = false;
let rootEl: HTMLElement;

const statusLabels: Record<string, string> = {
  new: 'Новая', assigned: 'Назначена', in_progress: 'В пути',
  completed: 'Завершена', cancelled: 'Отменена',
};

export async function renderApp(root: HTMLElement) {
  rootEl = root;
  if (!getToken()) { renderLogin(); return; }
  try {
    currentUser = await api.me();
    renderShell();
  } catch {
    clearToken();
    renderLogin();
  }
}

function renderLogin() {
  rootEl.innerHTML = `
    <div class="auth-page">
      <div class="auth-card">
        <h1><i class="bi bi-truck"></i> АТП — Учёт услуг</h1>
        <p class="subtitle">Автоматизированный учёт перевозок и контроль топлива</p>
        <form id="login-form">
          <div class="mb-3">
            <label class="form-label">Логин</label>
            <input class="form-control" name="login" value="dispatcher" required />
          </div>
          <div class="mb-3">
            <label class="form-label">Пароль</label>
            <input class="form-control" type="password" name="password" value="disp123" required />
          </div>
          <button type="submit" class="btn btn-primary w-100">Войти</button>
        </form>
        <p class="mt-3 small text-secondary">Демо: admin/admin123, dispatcher/disp123, driver1/drv123</p>
        <div id="login-error" class="text-danger mt-2"></div>
      </div>
    </div>`;
  fixDivs(rootEl);

  rootEl.querySelector('#login-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const fd = new FormData(e.target as HTMLFormElement);
    const errEl = rootEl.querySelector('#login-error')!;
    try {
      const res = await api.login(fd.get('login') as string, fd.get('password') as string);
      setToken(res.token);
      await renderApp(rootEl);
    } catch (err: unknown) {
      errEl.textContent = wailsError(err);
    }
  });
}

function fixDivs(el: HTMLElement) {
  el.innerHTML = el.innerHTML.replace(/<\/?div\b/g, (m) => m.replace('div', 'div'));
}

function navItems(role: Role) {
  if (role === 'driver') return [
    { id: 'orders', label: 'Мои рейсы', icon: 'bi-signpost-2' },
    { id: 'shift', label: 'Закрыть смену', icon: 'bi-fuel-pump' },
  ];
  if (role === 'admin') return [
    { id: 'dashboard', label: 'Панель', icon: 'bi-speedometer2' },
    { id: 'orders', label: 'Заявки', icon: 'bi-clipboard-check' },
    { id: 'new-order', label: 'Новая заявка', icon: 'bi-plus-circle' },
    { id: 'schedule', label: 'График', icon: 'bi-calendar3' },
    { id: 'tariffs', label: 'Тарифы', icon: 'bi-currency-exchange' },
    { id: 'vehicles', label: 'Автопарк', icon: 'bi-truck' },
    { id: 'reports', label: 'Отчёты', icon: 'bi-bar-chart' },
  ];
  return [
    { id: 'dashboard', label: 'Панель', icon: 'bi-speedometer2' },
    { id: 'orders', label: 'Заявки', icon: 'bi-clipboard-check' },
    { id: 'new-order', label: 'Новая заявка', icon: 'bi-plus-circle' },
    { id: 'schedule', label: 'График', icon: 'bi-calendar3' },
  ];
}

function roleLabel(role: Role) {
  return { admin: 'Админ', dispatcher: 'Диспетчер', driver: 'Водитель' }[role];
}

function renderShell() {
  const role = currentUser!.role;
  const nav = navItems(role);
  if (role === 'driver' && currentView === 'dashboard') currentView = 'orders';

  const navHtml = nav.map(n =>
    `<button type="button" class="nav-link ${currentView === n.id ? 'active' : ''}" data-view="${n.id}">
      <i class="bi ${n.icon}"></i> ${n.label}</button>`
  ).join('');

  rootEl.innerHTML = `
    <div class="app-shell">
      <nav class="sidebar ${sidebarOpen ? 'open' : ''}" id="sidebar">
        <div class="sidebar-brand"><i class="bi bi-truck-front"></i> АТП Учёт</div>
        ${navHtml}
        <button type="button" class="nav-link mt-auto text-danger" id="logout-btn">
          <i class="bi bi-box-arrow-right"></i> Выход</button>
      </nav>
      <div style="flex:1" class="d-flex flex-column">
        <header class="mobile-header">
          <button type="button" class="btn btn-outline-light btn-sm" id="menu-toggle"><i class="bi bi-list"></i></button>
          <strong>АТП Учёт</strong>
          <span class="badge bg-primary">${roleLabel(role)}</span>
        </header>
        <main class="main-content" id="main-view"></main>
      </div>
    </div>`;
  fixDivs(rootEl);

  rootEl.querySelectorAll('[data-view]').forEach(btn => {
    btn.addEventListener('click', () => {
      currentView = (btn as HTMLElement).dataset.view!;
      sidebarOpen = false;
      renderShell();
    });
  });
  rootEl.querySelector('#logout-btn')?.addEventListener('click', async () => {
    await api.logout().catch(() => {});
    clearToken();
    renderLogin();
  });
  rootEl.querySelector('#menu-toggle')?.addEventListener('click', () => {
    sidebarOpen = !sidebarOpen;
    rootEl.querySelector('#sidebar')?.classList.toggle('open', sidebarOpen);
  });

  renderView(rootEl.querySelector('#main-view') as HTMLElement, role);
}

async function renderView(main: HTMLElement, role: Role) {
  main.innerHTML = '<div class="text-center py-5"><div class="spinner-border text-primary"></div></div>';
  fixDivs(main);

  try {
    let html = '';
    switch (currentView) {
      case 'dashboard': html = await viewDashboard(); break;
      case 'orders': html = role === 'driver' ? await viewDriverOrders() : await viewOrders(); break;
      case 'new-order': html = await viewNewOrder(); break;
      case 'schedule': html = await viewSchedule(); break;
      case 'shift': html = await viewCloseShift(); break;
      case 'tariffs': html = await viewTariffs(); break;
      case 'vehicles': html = await viewVehicles(); break;
      case 'reports': html = await viewReports(); break;
      default: html = '<p>Раздел не найден</p>';
    }
    main.innerHTML = html;
    bindViewHandlers(main);
  } catch (err: unknown) {
    main.innerHTML = `<div class="alert alert-danger">${err instanceof Error ? err.message : 'Ошибка'}</div>`;
    fixDivs(main);
  }
}

function topBar(title: string) {
  return `<div class="top-bar"><div><h2 class="h4 mb-0">${title}</h2>
    <small class="text-secondary">${currentUser!.lastName} ${currentUser!.firstName}</small></div></div>`;
}

function formatMoney(n: number) {
  return new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(n);
}

async function viewDashboard() {
  const s = await api.dashboard();
  return topBar('Панель') + `<div class="stat-grid">
    <div class="stat-card"><div class="value">${s.ordersToday}</div><div class="label">Заявок сегодня</div></div>
    <div class="stat-card"><div class="value">${s.activeVehicles}</div><div class="label">Активных ТС</div></div>
    <div class="stat-card"><div class="value">${formatMoney(s.revenueMonth)}</div><div class="label">Выручка</div></div>
    <div class="stat-card"><div class="value text-warning">${s.fuelOverruns}</div><div class="label">Перерасходов</div></div>
  </div>`.replace(/<\/?div\b/g, m => m.replace('div', 'div')).replace(/div/g, 'div');
}

async function viewOrders() {
  const orders = await api.listOrders();
  const rows = orders.length ? orders.map(o => `
    <div class="card-panel order-card ${o.urgent ? 'urgent' : ''}">
      <div class="d-flex justify-content-between"><strong>${o.fromAddr} → ${o.toAddr}</strong>
      <span class="badge bg-secondary">${statusLabels[o.status] || o.status}</span></div>
      <div class="small text-secondary mt-1">${o.distanceKm} км · ${formatMoney(o.price)}</div>
    </div>`).join('') : '<p class="text-secondary">Нет заявок</p>';
  return topBar('Заявки') + rows.replace(/<\/?div\b/g, m => m.replace('div', 'div'));
}

async function viewDriverOrders() {
  const orders = await api.listOrders();
  const cards = orders.map(o => `
    <div class="card-panel order-card">
      <h5>${o.fromAddr}</h5><p><i class="bi bi-arrow-down"></i> ${o.toAddr}</p>
      <p class="small text-secondary">${new Date(o.scheduledAt).toLocaleString('ru-RU')}</p>
      ${o.status === 'assigned' ? `<button type="button" class="btn btn-success btn-trip" data-start="${o.id}">Начать рейс</button>` : ''}
      ${o.status === 'in_progress' ? `<button type="button" class="btn btn-primary btn-trip" data-complete="${o.id}">Завершить рейс</button>` : ''}
    </div>`).join('');
  return (topBar('Мои рейсы') + (cards || '<p>Нет рейсов</p>') +
    `<div class="driver-bottom-bar"><button type="button" class="btn btn-warning btn-trip" id="go-shift">Закрыть смену</button></div>`)
    .replace(/<\/?div\b/g, m => m.replace('div', 'div')).replace(/div/g, 'div');
}

function opts(items: Record<string, unknown>[], idKey: string, label: string | ((x: Record<string, unknown>) => string)) {
  return items.map(i => {
    const labelText = typeof label === 'function' ? label(i) : String(i[label]);
    return `<option value="${i[idKey]}">${labelText}</option>`;
  }).join('');
}

async function viewNewOrder() {
  const [clients, vehicles, tariffs, users] = await Promise.all([
    api.listClients(), api.listVehicles(), api.listTariffs(),
    api.listUsers().catch(() => [] as User[]),
  ]);
  const drivers = users.filter(u => u.role === 'driver');
  return (topBar('Новая заявка') + `
  <form id="order-form" class="card-panel"><div class="row g-3">
    <div class="col-md-6"><label class="form-label">Клиент</label>
      <select class="form-select" name="clientId" required>${opts(clients as any,'id','name')}</select></div>
    <div class="col-md-6"><label class="form-label">Тариф</label>
      <select class="form-select" name="tariffId" required>${opts(tariffs as any,'id','name')}</select></div>
    <div class="col-12"><label class="form-label">Откуда</label><input class="form-control" name="fromAddr" required /></div>
    <div class="col-12"><label class="form-label">Куда</label><input class="form-control" name="toAddr" required /></div>
    <div class="col-md-4"><label class="form-label">Км</label><input class="form-control" type="number" name="distanceKm" value="10" required /></div>
    <div class="col-md-4"><label class="form-label">Простой, ч</label><input class="form-control" type="number" name="idleHours" value="0" /></div>
    <div class="col-md-4"><label class="form-label">Срочная</label>
      <select class="form-select" name="urgent"><option value="false">Нет</option><option value="true">Да</option></select></div>
    <div class="col-md-6"><label class="form-label">Авто</label>
      <select class="form-select" name="vehicleId" required>${opts(vehicles.filter(v=>v.active) as any,'id','plate')}</select></div>
    <div class="col-md-6"><label class="form-label">Водитель</label>
      <select class="form-select" name="driverId" required>${opts(drivers as any,'id',u=>`${u.lastName} ${u.firstName}`)}</select></div>
    <div class="col-12"><div id="price-preview" class="alert alert-info">Цена: —</div></div>
    <div class="col-12"><button type="submit" class="btn btn-primary">Создать</button></div>
  </div></form>`).replace(/<\/?div\b/g, m => m.replace('div', 'div'));
}

async function viewSchedule() {
  const schedule = await api.vehicleSchedule();
  return (topBar('График') + schedule.map(s => `
    <div class="card-panel"><strong>${s.plate}</strong>
    ${s.orders.length ? s.orders.map(o => `<div class="schedule-slot">${o.fromAddr} → ${o.toAddr}</div>`).join('') : '<p class="text-secondary small">Свободен</p>'}
    </div>`).join('')).replace(/<\/?div\b/g, m => m.replace('div', 'div')).replace(/div/g, 'div');
}

async function viewCloseShift() {
  const vehicles = await api.listVehicles();
  return (topBar('Закрытие смены') + `
  <form id="shift-form" class="card-panel">
    <p class="small text-secondary">При перерасходе &gt;5% нужен комментарий.</p>
    <div class="row g-3">
      <div class="col-12"><label class="form-label">Авто</label>
        <select class="form-select" name="vehicleId" required>${opts(vehicles.filter(v=>v.active) as any,'id','plate')}</select></div>
      <div class="col-6"><label class="form-label">Пробег нач.</label><input class="form-control" type="number" name="startOdometer" required /></div>
      <div class="col-6"><label class="form-label">Пробег кон.</label><input class="form-control" type="number" name="endOdometer" required /></div>
      <div class="col-4"><label class="form-label">Топливо нач.</label><input class="form-control" type="number" name="fuelStart" required /></div>
      <div class="col-4"><label class="form-label">Заправка</label><input class="form-control" type="number" name="fuelRefilled" value="0" /></div>
      <div class="col-4"><label class="form-label">Топливо кон.</label><input class="form-control" type="number" name="fuelEnd" required /></div>
      <div class="col-12"><label class="form-label">Комментарий</label><textarea class="form-control" name="comment"></textarea></div>
      <div class="col-12"><button type="submit" class="btn btn-warning btn-trip">Закрыть смену</button></div>
    </div><div id="shift-result" class="mt-3"></div>
  </form>`).replace(/<\/?div\b/g, m => m.replace('div', 'div'));
}

async function viewTariffs() {
  const tariffs = await api.listTariffs();
  return topBar('Тарифы') + tariffs.map(t => `
    <div class="card-panel"><strong>${t.name}</strong>
    <p class="small text-secondary mb-0">Подача ${t.baseFee} ₽ · ${t.pricePerKm} ₽/км</p></div>`).join('').replace(/<\/?div\b/g, m => m.replace('div', 'div')) + `
  <form id="tariff-form" class="card-panel"><h5>Новый тариф</h5><div class="row g-2">
    <div class="col-md-4"><input class="form-control" name="name" placeholder="Название" required /></div>
    <div class="col-md-2"><input class="form-control" type="number" name="baseFee" placeholder="Подача" required /></div>
    <div class="col-md-2"><input class="form-control" type="number" name="pricePerKm" placeholder="₽/км" required /></div>
    <div class="col-md-2"><input class="form-control" type="number" name="pricePerIdleHr" placeholder="₽/ч" /></div>
    <div class="col-md-2"><button type="submit" class="btn btn-primary w-100">Сохранить</button></div>
  </div></form>`.replace(/<\/?div\b/g, m => m.replace('div', 'div'));
}

async function viewVehicles() {
  const vehicles = await api.listVehicles();
  return topBar('Автопарк') + vehicles.map(v =>
    `<div class="card-panel"><strong>${v.plate}</strong> — ${v.model}, ${v.fuelNorm} л/100км</div>`).join('').replace(/<\/?div\b/g, m => m.replace('div', 'div')).replace(/div/g,'div') + `
  <form id="vehicle-form" class="card-panel"><div class="row g-2">
    <div class="col-md-4"><input class="form-control" name="plate" placeholder="Номер" required /></div>
    <div class="col-md-4"><input class="form-control" name="model" placeholder="Модель" required /></div>
    <div class="col-md-2"><input class="form-control" type="number" name="fuelNorm" placeholder="л/100" required /></div>
    <div class="col-md-2"><button type="submit" class="btn btn-primary w-100">Добавить</button></div>
  </div></form>`.replace(/<\/?div\b/g, m => m.replace('div', 'div'));
}

async function viewReports() {
  const [stats, drivers] = await Promise.all([api.dashboard(), api.driverRating()]);
  const rows = drivers.map(d =>
    `<tr><td>${d.name}</td><td>${d.completed}/${d.total}</td><td>${Number(d.completionRate).toFixed(0)}%</td></tr>`).join('');
  return (topBar('Отчёты') + `
    <div class="stat-grid">
      <div class="stat-card"><div class="value">${formatMoney(stats.revenueMonth)}</div><div class="label">Выручка</div></div>
      <div class="stat-card"><div class="value">${stats.fuelOverruns}</div><div class="label">Перерасходов</div></div>
    </div>
    <div class="card-panel"><h5>Рейтинг водителей</h5>
    <table class="table table-sm"><thead><tr><th>Водитель</th><th>Рейсы</th><th>%</th></tr></thead>
    <tbody>${rows}</tbody></table></div>`).replace(/<\/?div\b/g, m => m.replace('div', 'div'));
}

function bindViewHandlers(main: HTMLElement) {
  main.querySelector('#go-shift')?.addEventListener('click', () => { currentView = 'shift'; renderShell(); });
  main.querySelectorAll('[data-start]').forEach(btn => btn.addEventListener('click', async () => {
    await api.updateOrderStatus((btn as HTMLElement).dataset.start!, 'in_progress'); renderShell();
  }));
  main.querySelectorAll('[data-complete]').forEach(btn => btn.addEventListener('click', async () => {
    await api.updateOrderStatus((btn as HTMLElement).dataset.complete!, 'completed'); renderShell();
  }));

  const orderForm = main.querySelector('#order-form');
  if (orderForm) {
    const updatePrice = async () => {
      const fd = new FormData(orderForm as HTMLFormElement);
      try {
        const price = await api.previewPrice(fd.get('tariffId') as string, Number(fd.get('distanceKm')),
          Number(fd.get('idleHours')), fd.get('urgent') === 'true');
        main.querySelector('#price-preview')!.textContent = `Цена: ${formatMoney(price)}`;
      } catch { main.querySelector('#price-preview')!.textContent = 'Цена: —'; }
    };
    orderForm.addEventListener('input', updatePrice);
    updatePrice();
    orderForm.addEventListener('submit', async e => {
      e.preventDefault();
      const fd = new FormData(e.target as HTMLFormElement);
      await api.createOrder({
        clientId: fd.get('clientId'), vehicleId: fd.get('vehicleId'), driverId: fd.get('driverId'),
        fromAddr: fd.get('fromAddr'), toAddr: fd.get('toAddr'),
        distanceKm: Number(fd.get('distanceKm')), idleHours: Number(fd.get('idleHours')),
        urgent: fd.get('urgent') === 'true', tariffId: fd.get('tariffId'),
        scheduledAt: new Date().toISOString(),
      });
      currentView = 'orders'; renderShell();
    });
  }

  main.querySelector('#shift-form')?.addEventListener('submit', async e => {
    e.preventDefault();
    const fd = new FormData(e.target as HTMLFormElement);
    const result = await api.closeShift({
      vehicleId: fd.get('vehicleId'), startOdometer: Number(fd.get('startOdometer')),
      endOdometer: Number(fd.get('endOdometer')), fuelStart: Number(fd.get('fuelStart')),
      fuelEnd: Number(fd.get('fuelEnd')), fuelRefilled: Number(fd.get('fuelRefilled')),
      comment: fd.get('comment'),
    });
    const el = main.querySelector('#shift-result')!;
    if (result.requireComment) el.innerHTML = `<div class="alert-fuel">${result.message}</div>`;
    else el.innerHTML = `<div class="alert alert-success">${result.message}</div>`;
  });

  main.querySelector('#tariff-form')?.addEventListener('submit', async e => {
    e.preventDefault();
    const fd = new FormData(e.target as HTMLFormElement);
    await api.saveTariff({ name: fd.get('name') as string, baseFee: Number(fd.get('baseFee')),
      pricePerKm: Number(fd.get('pricePerKm')), pricePerIdleHr: Number(fd.get('pricePerIdleHr')) || 0,
      urgencyCoeff: 1.5, active: true });
    currentView = 'tariffs'; renderShell();
  });

  main.querySelector('#vehicle-form')?.addEventListener('submit', async e => {
    e.preventDefault();
    const fd = new FormData(e.target as HTMLFormElement);
    await api.saveVehicle({ plate: fd.get('plate') as string, model: fd.get('model') as string,
      fuelNorm: Number(fd.get('fuelNorm')), active: true });
    currentView = 'vehicles'; renderShell();
  });
}
