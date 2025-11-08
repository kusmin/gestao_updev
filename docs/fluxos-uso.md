# Fluxos de Uso – Plataforma de Gestão Local

## Objetivo
Mapear jornadas essenciais do MVP para alinhar produto, UX e backend. Cada fluxo descreve atores, etapas, dados trocados e regras de validação.

## Fluxo 1 – Onboarding da Empresa
- **Ator**: Proprietário.
- **Passos**:
  1. Acessa `/signup` e informa dados da empresa (nome, CNPJ/CPF) + usuário admin (nome, e-mail, senha, telefone).
  2. Sistema envia e-mail de verificação (opcional no MVP) e cria tenant + usuário com role `admin`.
  3. Admin configura horários de funcionamento e serviços básicos (wizard em 3 etapas).
  4. Admin convida colaboradores informando e-mails ou gera link de convite.
- **Regras**:
  - Validar documento (CNPJ/CPF) e impedir duplicidade.
  - Senha mínima 8 caracteres com combinação de tipos.
  - Wizard pode ser pulado, mas dashboard exibe checklist até completar.

## Fluxo 2 – Agendamento de Serviço
- **Atores**: Recepcionista (manager/operator) ou cliente (futuro portal).
- **Passos**:
  1. Usuário seleciona profissional, serviço e data no calendário (visual semanal).
  2. Sistema carrega disponibilidade (bloqueios, folgas) e valida conflito.
  3. Usuário escolhe cliente existente ou cria registro rápido (nome, telefone).
  4. Define status inicial (`pending` ou `confirmed`) e adiciona notas.
  5. Sistema grava `booking` e opcionalmente dispara notificação (e-mail/WhatsApp manual).
- **Regras**:
  - Bloqueios de horário devem impedir agendamentos.
  - Profissional pode ter capacidade >1 (ex.: vendedor atende múltiplos clientes) via campo `max_parallel`.
  - Cancelamentos exigem motivo e podem gerar vagas automáticas.

## Fluxo 3 – Registro de Venda no Balcão
- **Ator**: Vendedor/Operador.
- **Passos**:
  1. Seleciona cliente (ou anonimiza) e adiciona itens: serviços concluídos + produtos.
  2. Sistema sugere itens do agendamento mais recente (se houver `booking_id`).
  3. Aplica descontos/vale fidelidade e calcula total.
  4. Escolhe método de pagamento (espécie, Pix manual, cartão) e registra pagamento.
  5. Sistema atualiza estoque (`inventory_movements`) e status `paid`.
- **Regras**:
  - Descontos >X% exigem permissão `manager`.
  - Estoque não pode ficar negativo; alertar se atingir nível crítico.
  - Pagamentos parcelados registrar múltiplas entradas em `payments`.

## Fluxo 4 – Gestão de Estoque
- **Ator**: Gerente.
- **Passos**:
  1. Acessa módulo de produtos e visualiza estoque atual + alertas.
  2. Registra entrada (compra fornecedor) ou ajuste (perda/dano) especificando motivo.
  3. Sistema cria `inventory_movement` e recalcula `stock_qty`.
  4. Pode exportar relatório (CSV) com histórico filtrado por período/produto.
- **Regras**:
  - Movimentos devem exigir justificativa textual.
  - Ajustes negativos grandes pedem dupla confirmação.
  - Exportações limitadas por `tenant_id` e respeitam timezone da empresa.

## Fluxo 5 – Relatórios Diários
- **Atores**: Admin/Manager.
- **Passos**:
  1. Ao abrir o dashboard, sistema agrega dados do dia: agendamentos realizados, atendimentos concluídos, vendas e recebimentos.
  2. Usuário aplica filtros (profissional, serviço, forma de pagamento).
  3. Sistema consulta views materializadas (`vw_daily_sales`, `vw_inventory_status`) e retorna cards + gráficos.
  4. Usuário pode enviar resumo por e-mail com um clique.
- **Regras**:
  - Acesso restrito a papéis `manager` e `admin`.
  - Dados devem ser atualizados em até 5 min (job de atualização ou consultas em tempo real).
  - Exportação de PDF/CSV programada para fases futuras.

## Notas de UX Geral
- Interface responsiva (desktop/tablet) com principais ações em até 3 cliques.
- Feedback imediato para erros de validação; mensagens internacionais via i18n (pt-BR default).
- Modo offline não previsto no MVP, mas cachê leve para listas recentes no frontend melhora percepção.

## Próximos Passos
1. Criar wireframes das telas críticas (signup, agenda, PDV, estoque, dashboard) e anexar imagens ou links no `docs/`.
2. Mapear eventos de analytics (ex.: `booking_created`, `sale_paid`) para instrumentação futura.
3. Elaborar fluxos adicionais: convites de usuários, lembretes automáticos, reativação de clientes inativos.
