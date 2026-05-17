import './style.css';
import { renderApp } from './ui/views';

document.addEventListener('DOMContentLoaded', () => {
  const app = document.querySelector('#app');
  if (app) renderApp(app as HTMLElement);
});
