import { api, clearToken, getToken, setToken, wailsError } from '../api/client';
import type { Client, Role, User } from '../types';

let currentUser: User | null = null;
let currentView = 'dashboard';
let lastCreatedClientId = '';
let rootEl: HTMLElement;

const statusLabels: Record<string, string> = {
  new: 'Новая',
  assigned: 'Назначена',
  in_progress: 'В пути',
  completed: 'Завершена',
  cancelled: 'Отменена',
};

const viewTitles: Record<string, string> = {
  dashboard: 'Обзор',
  orders: 'Рейсы',
  'new-order': 'Новая заявка',
  schedule: 'Загрузка парка',
  shift: 'Смена',
  tariffs: 'Тарифы',
  vehicles: 'Автопарк',
  clients: 'Клиенты',
  reports: 'Аналитика',
  users: 'Сотрудники',
};

export async function renderApp(root: HTMLElement) {
  rootEl = root;
  if (!getToken()) {
    renderLogin();
    return;
  }
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
    <div class="gate">
      <div class="gate-grid" aria-hidden="true"></div>
      <div class="gate-card">
        <div class="gate-logo">
          <div class="gate-logo-mark"><i class="bi bi-bezier2"></i></div>
          <div>
            <h1>TransitOS</h1>
            <span>dispatch · v2</span>
          </div>
        </div>
        <form id="login-form">
          <div class="mb-3">
            <label class="field-label">Оператор</label>
            <input class="form-control" name="login" value="dispatcher" required autocomplete="username" />
          </div>
          <div class="mb-4">
            <label class="field-label">Ключ доступа</label>
            <input class="form-control" type="password" name="password" value="disp123" required autocomplete="current-password" />
          </div>
          <button type="submit" class="btn btn-tx btn-tx-lg w-100">Подключиться</button>
        </form>
        <p class="gate-demo">admin / admin123<br/>dispatcher / disp123<br/>driver1 / drv123</p>
        <div id="login-error" class="text-danger mt-2 small"></div>
      </div>
    </div>`;

  rootEl.querySelector('#login-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const fd = new FormData(e.target as HTMLFormElement);
    const errEl = rootEl.querySelector('#login-error')!;
    errEl.textContent = '';
    try {
      const res = await api.login(fd.get('login') as string, fd.get('password') as string);
      setToken(res.token);
      await renderApp(rootEl);
    } catch (err: unknown) {
      errEl.textContent = wailsError(err);
    }
  });
}

function navItems(role: Role) {
  if (role === 'driver') {
    return [
      { id: 'orders', label: 'Рейсы', icon: 'bi-sign-turn-right' },
      { id: 'shift', label: 'Смена', icon: 'bi-speedometer2' },
    ];
  }
  if (role === 'admin') {
    return [
      { id: 'dashboard', label: 'Обзор', icon: 'bi-radar' },
      { id: 'orders', label: 'Рейсы', icon: 'bi-diagram-3' },
      { id: 'new-order', label: 'Создать', icon: 'bi-plus-lg', accent: true },
      { id: 'clients', label: 'Клиенты', icon: 'bi-building' },
      { id: 'schedule', label: 'Парк', icon: 'bi-columns-gap' },
      { id: 'tariffs', label: 'Тарифы', icon: 'bi-sliders' },
      { id: 'vehicles', label: 'ТС', icon: 'bi-bus-front' },
      { id: 'users', label: 'Сотрудники', icon: 'bi-people' },
      { id: 'reports', label: 'Отчёты', icon: 'bi-graph-up-arrow' },
    ];
  }
  return [
    { id: 'dashboard', label: 'Обзор', icon: 'bi-radar' },
    { id: 'orders', label: 'Рейсы', icon: 'bi-diagram-3' },
    { id: 'new-order', label: 'Создать', icon: 'bi-plus-lg', accent: true },
    { id: 'clients', label: 'Клиенты', icon: 'bi-building' },
    { id: 'schedule', label: 'Парк', icon: 'bi-columns-gap' },
  ];
}

function roleLabel(role: Role) {
  return { admin: 'Администратор', dispatcher: 'Диспетчер', driver: 'Водитель' }[role];
}

function roleChip(role: Role) {
  return `<span class="chip role-${role}">${roleLabel(role)}</span>`;
}

function userInitials() {
  const u = currentUser!;
  return `${u.firstName?.[0] || ''}${u.lastName?.[0] || ''}`.toUpperCase() || '?';
}

function dockTab(n: { id: string; label: string; icon: string; accent?: boolean }) {
  const active = currentView === n.id;
  const cls = ['dock-tab', active ? 'active' : '', n.accent ? 'accent-new' : ''].filter(Boolean).join(' ');
  return `<button type="button" class="${cls}" data-view="${n.id}"><i class="bi ${n.icon}"></i>${n.label}</button>`;
}

function renderShell() {
  const role = currentUser!.role;
  const isDriver = role === 'driver';
  if (isDriver && currentView === 'dashboard') currentView = 'orders';

  const nav = navItems(role);
  const dockHtml = nav.map(dockTab).join('');

  rootEl.innerHTML = `
    <div class="hub ${isDriver ? 'driver-mode' : ''}">
      <header class="hub-header">
        <div class="hub-brand">
          <i class="bi bi-bezier2"></i>
          <strong>TransitOS</strong>
        </div>
        <div class="hub-meta">
          <div class="hub-user">
            <div class="name">${currentUser!.firstName} ${currentUser!.lastName}</div>
            <div class="role">${roleLabel(role)}</div>
          </div>
          <div class="hub-avatar">${userInitials()}</div>
          <button type="button" class="btn-ghost" id="logout-btn" title="Выход"><i class="bi bi-power"></i></button>
        </div>
      </header>
      <nav class="hub-dock" id="hub-dock">${dockHtml}</nav>
      <nav class="hub-dock-bottom" id="hub-dock-bottom">${dockHtml}</nav>
      <main class="hub-stage" id="main-view"></main>
    </div>`;

  rootEl.querySelectorAll('[data-view]').forEach((btn) => {
    btn.addEventListener('click', () => {
      currentView = (btn as HTMLElement).dataset.view!;
      renderShell();
    });
  });
  rootEl.querySelector('#logout-btn')?.addEventListener('click', async () => {
    await api.logout().catch(() => {});
    clearToken();
    renderLogin();
  });

  renderView(rootEl.querySelector('#main-view') as HTMLElement, role);
}

function stageHeader(subtitle?: string) {
  return `
    <p class="stage-title">${subtitle || 'модуль'}</p>
    <h2 class="stage-heading">${viewTitles[currentView] || 'Раздел'}</h2>`;
}

async function renderView(main: HTMLElement, role: Role) {
  main.innerHTML = `<div class="loader-pulse"><span></span><span></span><span></span> синхронизация</div>`;

  try {
    let html = '';
    switch (currentView) {
      case 'dashboard':
        html = await viewDashboard();
        break;
      case 'orders':
        html = role === 'driver' ? await viewDriverOrders() : await viewOrders();
        break;
      case 'new-order':
        html = await viewNewOrder();
        break;
      case 'schedule':
        html = await viewSchedule();
        break;
      case 'shift':
        html = await viewCloseShift();
        break;
      case 'tariffs':
        html = await viewTariffs();
        break;
      case 'vehicles':
        html = await viewVehicles();
        break;
      case 'clients':
        html = await viewClients();
        break;
      case 'reports':
        html = await viewReports();
        break;
      case 'users':
        html = await viewUsers();
        break;
      default:
        html = stageHeader() + '<div class="empty"><i class="bi bi-question-lg"></i><p>Раздел не найден</p></div>';
    }
    main.innerHTML = html;
    bindViewHandlers(main);
  } catch (err: unknown) {
    main.innerHTML = stageHeader('ошибка') + `<div class="alert-tx-warn">${wailsError(err)}</div>`;
  }
}

function formatMoney(n: number) {
  return new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(n);
}

function statusChip(status: string, urgent?: boolean) {
  const chips = [urgent ? '<span class="chip urgent">срочно</span>' : ''];
  chips.push(`<span class="chip status-${status}">${statusLabels[status] || status}</span>`);
  return chips.join('');
}

async function viewDashboard() {
  const s = await api.dashboard();
  return (
    stageHeader('сводка за сутки') +
    `
    <div class="bento">
      <div class="bento-cell bento-hero">
        <div class="lbl">Заявок сегодня</div>
        <div class="kpi">${s.ordersToday}</div>
      </div>
      <div class="bento-cell">
        <div class="lbl">Активных ТС</div>
        <div class="kpi-sm">${s.activeVehicles}</div>
      </div>
      <div class="bento-cell">
        <div class="lbl">Выручка / мес</div>
        <div class="kpi-sm">${formatMoney(s.revenueMonth)}</div>
      </div>
      <div class="bento-cell warn">
        <div class="lbl">Перерасход топлива</div>
        <div class="kpi-sm">${s.fuelOverruns}</div>
      </div>
    </div>`
  );
}

async function viewOrders() {
  const orders = (await api.listOrders()) ?? [];
  if (!orders.length) {
    return stageHeader('лента рейсов') + '<div class="empty"><i class="bi bi-inbox"></i><p>Нет активных заявок</p></div>';
  }
  const feed = orders
    .map(
      (o) => `
      <article class="route-card ${o.urgent ? 'urgent' : ''} ${o.status}">
        <div class="route-head">
          <div class="route-path">${o.fromAddr}<span class="arrow">→</span>${o.toAddr}</div>
          ${statusChip(o.status, o.urgent)}
        </div>
        <div class="route-meta">${o.distanceKm} км · ${formatMoney(o.price)}</div>
      </article>`
    )
    .join('');
  return stageHeader('лента рейсов') + `<div class="route-feed">${feed}</div>`;
}

async function viewDriverOrders() {
  const [ordersRaw, shift] = await Promise.all([api.listOrders(), api.shiftStatus()]);
  const orders = ordersRaw ?? [];
  const shiftClosed = shift.closed;
  const shiftOpened = shift.opened && !shift.closed;
  const cards = orders
    .map(
      (o) => `
      <div class="trip-mobile">
        <div class="from">${o.fromAddr}</div>
        <div class="to-line"><i class="bi bi-arrow-down-short"></i> ${o.toAddr}</div>
        <div class="route-meta mb-3">${o.scheduledAt ? new Date(o.scheduledAt).toLocaleString('ru-RU') : '—'}</div>
        ${o.status === 'assigned' ? `<button type="button" class="btn btn-tx-success" data-start="${o.id}">Старт рейса</button>` : ''}
        ${o.status === 'in_progress' ? `<button type="button" class="btn btn-tx" data-complete="${o.id}">Завершить</button>` : ''}
      </div>`
    )
    .join('');

  let shiftFab = '';
  if (shiftClosed) {
    shiftFab = '<div class="driver-fab-wrap"><div class="alert-tx-ok">Смена закрыта на сегодня</div></div>';
  } else if (!shiftOpened) {
    shiftFab =
      '<div class="driver-fab-wrap"><button type="button" class="btn btn-tx btn-tx-lg" id="go-shift-open">Открыть смену</button></div>';
  } else {
    shiftFab =
      '<div class="driver-fab-wrap"><button type="button" class="btn btn-tx btn-tx-lg" id="go-shift-close">Закрыть смену</button></div>';
  }

  const body =
    (cards || '<div class="empty"><i class="bi bi-calendar-x"></i><p>Нет рейсов на сегодня</p></div>') + shiftFab;
  return stageHeader('ваши рейсы') + `<div class="trip-stack">${body}</div>`;
}

function opts(
  items: { id?: string }[],
  idKey: string,
  label: string | ((x: Record<string, unknown>) => string),
  selectedId = ''
) {
  return items
    .map((i) => {
      const rec = i as Record<string, unknown>;
      const id = String(rec[idKey] ?? rec.ID ?? rec.id ?? '');
      const labelText = typeof label === 'function' ? label(rec) : String(rec[label]);
      const sel = selectedId && id === selectedId ? ' selected' : '';
      return `<option value="${id}"${sel}>${labelText}</option>`;
    })
    .join('');
}

function clientOpts(clients: Client[]) {
  if (!clients.length) {
    return '<option value="" disabled selected>Нет клиентов — добавьте ниже</option>';
  }
  return opts(clients, 'id', 'name', lastCreatedClientId);
}

function clientFormFields(msgId: string) {
  return `
    <div class="sheet-grid">
      <div class="form-field" style="grid-column: 1 / -1">
        <label>Название / организация</label>
        <input class="form-control" name="name" required placeholder="ООО «Логистика»" />
      </div>
      <div class="form-field">
        <label>Телефон</label>
        <input class="form-control" name="phone" type="tel" placeholder="+7..." />
      </div>
      <div class="form-field">
        <label>Лимит долга, ₽</label>
        <input class="form-control" type="number" name="debtLimit" value="0" min="0" step="1000" />
      </div>
      <div class="form-field" style="grid-column: 1 / -1">
        <div id="${msgId}" class="mt-1"></div>
        <button type="submit" class="btn btn-tx">Сохранить клиента</button>
      </div>
    </div>`;
}

async function viewNewOrder() {
  const [clients, vehicles, tariffs, drivers] = await Promise.all([
    api.listClients(),
    api.listVehicles(),
    api.listTariffs(),
    api.listDriversAvailable().catch(() => [] as User[]),
  ]);

  return (
    stageHeader('мастер заявки') +
    `
    <div class="sheet">
      <p class="sheet-title">Параметры перевозки</p>
      <form id="order-form">
        <div class="sheet-grid">
          <div class="form-field">
            <label>Клиент</label>
            <select class="form-select" name="clientId" required>${clientOpts(clients)}</select>
          </div>
          <div class="form-field">
            <label>Тариф</label>
            <select class="form-select" name="tariffId" required>${opts(tariffs, 'id', 'name')}</select>
          </div>
          <div class="form-field" style="grid-column: 1 / -1">
            <label>Откуда</label>
            <input class="form-control" name="fromAddr" required placeholder="Адрес погрузки" />
          </div>
          <div class="form-field" style="grid-column: 1 / -1">
            <label>Куда</label>
            <input class="form-control" name="toAddr" required placeholder="Адрес выгрузки" />
          </div>
          <div class="form-field">
            <label>Км</label>
            <input class="form-control" type="number" name="distanceKm" value="10" required />
          </div>
          <div class="form-field">
            <label>Простой, ч</label>
            <input class="form-control" type="number" name="idleHours" value="0" step="0.1" />
          </div>
          <div class="form-field">
            <label>Срочная</label>
            <select class="form-select" name="urgent">
              <option value="false">Нет</option>
              <option value="true">Да</option>
            </select>
          </div>
          <div class="form-field">
            <label>Автомобиль</label>
            <select class="form-select" name="vehicleId" required>${opts(vehicles.filter((v) => v.active), 'id', 'plate')}</select>
          </div>
          <div class="form-field">
            <label>Водитель</label>
            <select class="form-select" name="driverId" required>${opts(drivers, 'id', (u) => `${u.lastName} ${u.firstName}`)}</select>
          </div>
          <div class="form-field" style="grid-column: 1 / -1">
            <div id="price-preview" class="price-box">Расчёт: —</div>
          </div>
        </div>
        <button type="submit" class="btn btn-tx mt-3">Создать заявку</button>
      </form>
    </div>
    <div class="sheet mt-3">
      <p class="sheet-title">Новый клиент</p>
      <form id="client-quick-form">${clientFormFields('client-quick-msg')}</form>
    </div>`
  );
}

async function viewSchedule() {
  const schedule = await api.vehicleSchedule();
  if (!schedule.length) {
    return stageHeader('таймлайн') + '<div class="empty"><i class="bi bi-truck"></i><p>Нет активных ТС</p></div>';
  }

  const rows = schedule
    .map((item) => {
      const orders = item.orders ?? [];
      const blocks =
        orders.length > 0
          ? orders
              .map(
                (o) =>
                  `<div class="gantt-block"><strong>${o.fromAddr}</strong>${o.toAddr}<br/><span class="chip status-${o.status}">${statusLabels[o.status] || o.status}</span></div>`
              )
              .join('')
          : '<div class="gantt-block empty">свободен</div>';
      return `
        <div class="gantt-row">
          <div class="gantt-plate">${item.plate}</div>
          <div class="gantt-track">${blocks}</div>
        </div>`;
    })
    .join('');

  return stageHeader('таймлайн') + `<div class="gantt-board">${rows}</div>`;
}

async function viewCloseShift() {
  const [vehicles, shift] = await Promise.all([api.listVehicles(), api.shiftStatus()]);
  if (shift.closed) {
    return (
      stageHeader('учёт смены') +
      `<div class="sheet"><div class="alert-tx-ok">Смена закрыта на ${shift.date}. Новые рейсы сегодня назначаться не будут.</div></div>`
    );
  }
  if (!shift.opened) {
    return (
      stageHeader('открытие смены') +
      `
      <div class="sheet">
        <p class="sheet-title">Начало смены · укажите автомобиль и показания одометра</p>
        <form id="shift-open-form">
          <div class="sheet-grid">
            <div class="form-field" style="grid-column: 1 / -1">
              <label>Автомобиль</label>
              <select class="form-select" name="vehicleId" required>${opts(vehicles.filter((v) => v.active), 'id', 'plate')}</select>
            </div>
            <div class="form-field">
              <label>Пробег нач., км</label>
              <input class="form-control" type="number" name="startOdometer" required />
            </div>
            <div class="form-field">
              <label>Топливо в баке, л</label>
              <input class="form-control" type="number" step="0.1" name="fuelStart" required />
            </div>
          </div>
          <button type="submit" class="btn btn-tx btn-tx-lg mt-3">Открыть смену</button>
          <div id="shift-result" class="mt-3"></div>
        </form>
      </div>`
    );
  }
  const startKm = shift.startOdometer ?? 0;
  const startFuel = shift.fuelStart ?? 0;
  return (
    stageHeader('закрытие смены') +
    `
    <div class="sheet mb-3">
      <p class="sheet-title">Смена открыта</p>
      <p class="route-meta">Пробег нач.: ${startKm} км · Топливо нач.: ${startFuel} л</p>
    </div>
    <div class="sheet">
      <p class="sheet-title">Конец смены · при перерасходе &gt;5% нужен комментарий</p>
      <form id="shift-form">
        <div class="sheet-grid">
          <div class="form-field">
            <label>Пробег кон., км</label>
            <input class="form-control" type="number" name="endOdometer" min="${startKm + 1}" required />
          </div>
          <div class="form-field">
            <label>Заправка, л</label>
            <input class="form-control" type="number" step="0.1" name="fuelRefilled" value="0" />
          </div>
          <div class="form-field">
            <label>Топливо кон., л</label>
            <input class="form-control" type="number" step="0.1" name="fuelEnd" required />
          </div>
          <div class="form-field" style="grid-column: 1 / -1">
            <label>Комментарий</label>
            <textarea class="form-control" name="comment" rows="2"></textarea>
          </div>
        </div>
        <button type="submit" class="btn btn-tx btn-tx-lg mt-3">Закрыть смену</button>
        <div id="shift-result" class="mt-3"></div>
      </form>
    </div>`
  );
}

async function viewTariffs() {
  const tariffs = await api.listTariffs();
  const list = tariffs
    .map(
      (t) => `
      <div class="list-row">
        <div><strong>${t.name}</strong><div class="route-meta">подача ${t.baseFee} ₽ · ${t.pricePerKm} ₽/км · простой ${t.pricePerIdleHr} ₽/ч</div></div>
      </div>`
    )
    .join('');

  return (
    stageHeader('тарифная сетка') +
    `
    <div class="sheet mb-3">${list || '<div class="empty"><p>Нет тарифов</p></div>'}</div>
    <div class="sheet">
      <p class="sheet-title">Новый тариф</p>
      <form id="tariff-form" class="sheet-grid">
        <div class="form-field"><label>Название</label><input class="form-control" name="name" required /></div>
        <div class="form-field"><label>Подача</label><input class="form-control" type="number" name="baseFee" required /></div>
        <div class="form-field"><label>₽/км</label><input class="form-control" type="number" name="pricePerKm" required /></div>
        <div class="form-field"><label>₽/ч простой</label><input class="form-control" type="number" name="pricePerIdleHr" /></div>
        <div class="form-field" style="align-self:end"><button type="submit" class="btn btn-tx w-100">Сохранить</button></div>
      </form>
    </div>`
  );
}

async function viewVehicles() {
  const vehicles = await api.listVehicles();
  const list = vehicles
    .map(
      (v) => `
      <div class="list-row">
        <div><strong class="gantt-plate" style="display:inline;padding:0.2rem 0.5rem;margin-right:0.5rem">${v.plate}</strong> ${v.model}</div>
        <span class="chip">${v.fuelNorm} л/100</span>
      </div>`
    )
    .join('');

  return (
    stageHeader('реестр ТС') +
    `
    <div class="sheet mb-3">${list}</div>
    <div class="sheet">
      <p class="sheet-title">Добавить транспорт</p>
      <form id="vehicle-form" class="sheet-grid">
        <div class="form-field"><label>Госномер</label><input class="form-control" name="plate" required /></div>
        <div class="form-field"><label>Модель</label><input class="form-control" name="model" required /></div>
        <div class="form-field"><label>л/100км</label><input class="form-control" type="number" step="0.1" name="fuelNorm" required /></div>
        <div class="form-field" style="align-self:end"><button type="submit" class="btn btn-tx w-100">Добавить</button></div>
      </form>
    </div>`
  );
}

async function viewClients() {
  const clients = await api.listClients();
  const sorted = [...clients].sort((a, b) => a.name.localeCompare(b.name, 'ru'));
  const list = sorted
    .map(
      (c) => `
      <div class="list-row">
        <div>
          <strong>${c.name}</strong>
          <div class="route-meta">${c.phone || 'телефон не указан'} · лимит ${formatMoney(c.debtLimit)}</div>
        </div>
      </div>`
    )
    .join('');

  return (
    stageHeader('справочник') +
    `
    <div class="sheet mb-3">
      <p class="sheet-title">В системе · ${clients.length}</p>
      ${list || '<div class="empty"><i class="bi bi-building"></i><p>Нет клиентов</p></div>'}
    </div>
    <div class="sheet">
      <p class="sheet-title">Добавить клиента</p>
      <form id="client-form">${clientFormFields('client-form-msg')}</form>
    </div>`
  );
}

async function viewUsers() {
  const users = await api.listUsers();
  const sorted = [...users].sort((a, b) => {
    const order: Record<Role, number> = { admin: 0, dispatcher: 1, driver: 2 };
    const d = (order[a.role] ?? 9) - (order[b.role] ?? 9);
    return d !== 0 ? d : a.login.localeCompare(b.login);
  });

  const list = sorted
    .map((u) => {
      const isSelf = u.id === currentUser?.id;
      const deleteBtn = isSelf
        ? ''
        : `<button type="button" class="btn btn-ghost btn-delete-user" data-user-id="${u.id}" data-user-login="${u.login}" title="Удалить"><i class="bi bi-trash"></i></button>`;
      return `
      <div class="list-row">
        <div>
          <strong>${u.lastName} ${u.firstName}</strong>
          <div class="route-meta">@${u.login}${u.phone ? ` · ${u.phone}` : ''}</div>
        </div>
        <div class="list-row-actions">
          ${roleChip(u.role)}
          ${deleteBtn}
        </div>
      </div>`;
    })
    .join('');

  return (
    stageHeader('учётные записи') +
    `
    <div class="sheet mb-3">
      <p class="sheet-title">В системе · ${users.length}</p>
      ${list || '<div class="empty"><p>Нет пользователей</p></div>'}
    </div>
    <div class="sheet">
      <p class="sheet-title">Новый сотрудник</p>
      <form id="user-form">
        <div class="sheet-grid">
          <div class="form-field">
            <label>Роль</label>
            <select class="form-select" name="role" required>
              <option value="driver">Водитель</option>
              <option value="dispatcher">Диспетчер</option>
              <option value="admin">Администратор</option>
            </select>
          </div>
          <div class="form-field">
            <label>Логин</label>
            <input class="form-control" name="login" required autocomplete="off" placeholder="например driver3" />
          </div>
          <div class="form-field">
            <label>Пароль</label>
            <input class="form-control" type="password" name="password" required minlength="4" autocomplete="new-password" />
          </div>
          <div class="form-field">
            <label>Имя</label>
            <input class="form-control" name="firstName" required />
          </div>
          <div class="form-field">
            <label>Фамилия</label>
            <input class="form-control" name="lastName" required />
          </div>
          <div class="form-field">
            <label>Телефон</label>
            <input class="form-control" name="phone" type="tel" placeholder="+7..." />
          </div>
        </div>
        <div id="user-form-msg" class="mt-2"></div>
        <button type="submit" class="btn btn-tx mt-3">Создать учётную запись</button>
      </form>
    </div>`
  );
}

async function viewReports() {
  const [stats, drivers] = await Promise.all([api.dashboard(), api.driverRating()]);
  const rows = drivers
    .map(
      (d) =>
        `<tr><td>${d.name}</td><td>${d.completed} / ${d.total}</td><td><strong>${Number(d.completionRate).toFixed(0)}%</strong></td></tr>`
    )
    .join('');

  return (
    stageHeader('аналитика') +
    `
    <div class="bento" style="margin-bottom:1rem">
      <div class="bento-cell"><div class="lbl">Выручка</div><div class="kpi-sm">${formatMoney(stats.revenueMonth)}</div></div>
      <div class="bento-cell warn"><div class="lbl">Перерасходов</div><div class="kpi-sm">${stats.fuelOverruns}</div></div>
    </div>
    <div class="data-table-wrap">
      <table class="data-table">
        <thead><tr><th>Водитель</th><th>Рейсы</th><th>Выполнение</th></tr></thead>
        <tbody>${rows}</tbody>
      </table>
    </div>`
  );
}

function bindViewHandlers(main: HTMLElement) {
  main.querySelector('#go-shift-open')?.addEventListener('click', () => {
    currentView = 'shift';
    renderShell();
  });
  main.querySelector('#go-shift-close')?.addEventListener('click', () => {
    currentView = 'shift';
    renderShell();
  });

  main.querySelectorAll('[data-start]').forEach((btn) => {
    btn.addEventListener('click', async () => {
      await api.updateOrderStatus((btn as HTMLElement).dataset.start!, 'in_progress');
      renderShell();
    });
  });
  main.querySelectorAll('[data-complete]').forEach((btn) => {
    btn.addEventListener('click', async () => {
      await api.updateOrderStatus((btn as HTMLElement).dataset.complete!, 'completed');
      renderShell();
    });
  });

  const orderForm = main.querySelector('#order-form');
  if (orderForm) {
    const updatePrice = async () => {
      const fd = new FormData(orderForm as HTMLFormElement);
      const el = main.querySelector('#price-preview');
      if (!el) return;
      try {
        const price = await api.previewPrice(
          fd.get('tariffId') as string,
          Number(fd.get('distanceKm')),
          Number(fd.get('idleHours')),
          fd.get('urgent') === 'true'
        );
        el.textContent = `Расчёт: ${formatMoney(price)}`;
      } catch {
        el.textContent = 'Расчёт: недоступен';
      }
    };
    orderForm.addEventListener('input', updatePrice);
    updatePrice();
    orderForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const fd = new FormData(e.target as HTMLFormElement);
      await api.createOrder({
        clientId: fd.get('clientId'),
        vehicleId: fd.get('vehicleId'),
        driverId: fd.get('driverId'),
        fromAddr: fd.get('fromAddr'),
        toAddr: fd.get('toAddr'),
        distanceKm: Number(fd.get('distanceKm')),
        idleHours: Number(fd.get('idleHours')),
        urgent: fd.get('urgent') === 'true',
        tariffId: fd.get('tariffId'),
        scheduledAt: new Date().toISOString(),
      });
      currentView = 'orders';
      renderShell();
    });
  }


  main.querySelector('#shift-open-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const fd = new FormData(e.target as HTMLFormElement);
    const el = main.querySelector('#shift-result')!;
    el.innerHTML = '';
    try {
      const result = await api.openShift({
        vehicleId: String(fd.get('vehicleId') || ''),
        startOdometer: Number(fd.get('startOdometer')),
        fuelStart: Number(fd.get('fuelStart')),
      });
      el.innerHTML = `<div class="alert-tx-ok">${result.message}</div>`;
      setTimeout(() => renderShell(), 800);
    } catch (err: unknown) {
      el.innerHTML = `<div class="alert-tx-warn">${wailsError(err)}</div>`;
    }
  });

  main.querySelector('#shift-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const fd = new FormData(e.target as HTMLFormElement);
    const el = main.querySelector('#shift-result')!;
    el.innerHTML = '';
    try {
      const result = await api.closeShift({
        endOdometer: Number(fd.get('endOdometer')),
        fuelEnd: Number(fd.get('fuelEnd')),
        fuelRefilled: Number(fd.get('fuelRefilled') || 0),
        comment: String(fd.get('comment') || ''),
      });
      el.innerHTML = result.requireComment
        ? `<div class="alert-tx-warn">${result.message}</div>`
        : `<div class="alert-tx-ok">${result.message}</div>`;
      if (!result.requireComment) {
        setTimeout(() => renderShell(), 800);
      }
    } catch (err: unknown) {
      el.innerHTML = `<div class="alert-tx-warn">${wailsError(err)}</div>`;
    }
  });

  main.querySelector('#tariff-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const fd = new FormData(e.target as HTMLFormElement);
    await api.saveTariff({
      name: fd.get('name') as string,
      baseFee: Number(fd.get('baseFee')),
      pricePerKm: Number(fd.get('pricePerKm')),
      pricePerIdleHr: Number(fd.get('pricePerIdleHr')) || 0,
      urgencyCoeff: 1.5,
      active: true,
    });
    currentView = 'tariffs';
    renderShell();
  });

  main.querySelector('#vehicle-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const fd = new FormData(e.target as HTMLFormElement);
    await api.saveVehicle({
      plate: fd.get('plate') as string,
      model: fd.get('model') as string,
      fuelNorm: Number(fd.get('fuelNorm')),
      active: true,
    });
    currentView = 'vehicles';
    renderShell();
  });

  const bindClientForm = (formId: string, msgId: string, stayView: string) => {
    const form = main.querySelector(formId);
    if (!form) return;
    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      const fd = new FormData(e.target as HTMLFormElement);
      const msg = main.querySelector(msgId)!;
      msg.innerHTML = '';
      try {
        const saved = await api.saveClient({
          name: fd.get('name') as string,
          phone: (fd.get('phone') as string) || '',
          debtLimit: Number(fd.get('debtLimit')) || 0,
        });
        lastCreatedClientId = saved.id || '';
        msg.innerHTML = '<div class="alert-tx-ok">Клиент сохранён</div>';
        (e.target as HTMLFormElement).reset();
        currentView = stayView;
        setTimeout(() => renderShell(), 600);
      } catch (err: unknown) {
        msg.innerHTML = `<div class="alert-tx-warn">${wailsError(err)}</div>`;
      }
    });
  };
  bindClientForm('#client-form', '#client-form-msg', 'clients');
  bindClientForm('#client-quick-form', '#client-quick-msg', 'new-order');

  main.querySelectorAll('.btn-delete-user').forEach((btn) => {
    btn.addEventListener('click', async () => {
      const el = btn as HTMLElement;
      const id = el.dataset.userId!;
      const login = el.dataset.userLogin!;
      if (!confirm(`Удалить учётную запись @${login}? Это действие нельзя отменить.`)) return;
      try {
        await api.deleteUser(id);
        renderShell();
      } catch (err: unknown) {
        alert(wailsError(err));
      }
    });
  });

  main.querySelector('#user-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const fd = new FormData(e.target as HTMLFormElement);
    const msg = main.querySelector('#user-form-msg')!;
    msg.innerHTML = '';
    try {
      await api.createUser(
        {
          login: fd.get('login') as string,
          role: fd.get('role') as Role,
          firstName: fd.get('firstName') as string,
          lastName: fd.get('lastName') as string,
          phone: (fd.get('phone') as string) || '',
        },
        fd.get('password') as string
      );
      msg.innerHTML = '<div class="alert-tx-ok">Учётная запись создана</div>';
      (e.target as HTMLFormElement).reset();
      setTimeout(() => renderShell(), 800);
    } catch (err: unknown) {
      msg.innerHTML = `<div class="alert-tx-warn">${wailsError(err)}</div>`;
    }

  });
}

