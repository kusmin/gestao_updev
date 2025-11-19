import React from 'react';
import FullCalendar from '@fullcalendar/react';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';

function ExampleCalendar() {
  const handleDateClick = (arg: any) => {
    // bind with an arrow function
    alert(arg.dateStr);
  };

  return (
    <div className="p-4">
      <FullCalendar
        plugins={[dayGridPlugin, interactionPlugin]}
        initialView="dayGridMonth"
        weekends={true}
        events={[
          { title: 'evento 1', date: '2025-11-19' },
          { title: 'evento 2', date: '2025-11-20' },
        ]}
        dateClick={handleDateClick}
        eventClick={(info) => alert(`Evento: ${info.event.title}`)}
      />
    </div>
  );
}

export default ExampleCalendar;
