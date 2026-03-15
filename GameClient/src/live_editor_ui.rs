use egui::{CentralPanel, Color32, FontId, RichText, ScrollArea, TextEdit, TopBottomPanel};
use std::sync::{Arc, Mutex};
use tokio::sync::mpsc::UnboundedSender;

pub struct EditorState {
    pub equation:    String,
    pub log:         String,
    pub ws_tx:       Arc<Mutex<Option<UnboundedSender<String>>>>,
}

impl EditorState {
    pub fn new(ws_tx: Arc<Mutex<Option<UnboundedSender<String>>>>) -> Self {
        Self {
            equation: String::from("sin(x) - sqrt(y**2 + z**2) - 1.0"),
            log:      String::from("[ ax10m firewall ready ]\n"),
            ws_tx,
        }
    }

    pub fn draw(&mut self, ctx: &egui::Context) {
        TopBottomPanel::bottom("terminal_panel")
            .min_height(96.0)
            .show(ctx, |ui| {
                ui.label(RichText::new("─── Firewall Output ───").color(Color32::DARK_GRAY));
                ScrollArea::vertical()
                    .stick_to_bottom(true)
                    .show(ui, |ui| {
                        ui.add(
                            TextEdit::multiline(&mut self.log.as_str())
                                .font(FontId::monospace(11.0))
                                .text_color(Color32::LIGHT_GREEN)
                                .desired_rows(4)
                                .desired_width(f32::INFINITY),
                        );
                    });
            });

        CentralPanel::default().show(ctx, |ui| {
            ui.label(RichText::new("ax10m  ·  Live Math Editor").strong().size(15.0));
            ui.separator();

            let available = ui.available_height() - 38.0;
            ScrollArea::vertical()
                .max_height(available)
                .show(ui, |ui| {
                    ui.add(
                        TextEdit::multiline(&mut self.equation)
                            .font(FontId::monospace(13.0))
                            .desired_rows(8)
                            .desired_width(f32::INFINITY)
                            .hint_text("sin(x) + cos(y) - z"),
                    );
                });

            ui.separator();
            if ui.button("  Compile & Inject").clicked() {
                self.compile_and_send();
            }
        });
    }

    fn compile_and_send(&mut self) {
        let payload = serde_json::json!({ "equation": self.equation.trim() }).to_string();
        let guard = self.ws_tx.lock().unwrap();
        if let Some(tx) = guard.as_ref() {
            match tx.send(payload) {
                Ok(_)  => self.log.push_str("[TX] equation dispatched\n"),
                Err(e) => self.log.push_str(&format!("[ERR] send failed: {e}\n")),
            }
        } else {
            self.log.push_str("[WARN] WebSocket not connected\n");
        }
    }

    pub fn push_log(&mut self, msg: &str) {
        self.log.push_str(msg);
        self.log.push('\n');
    }
}
