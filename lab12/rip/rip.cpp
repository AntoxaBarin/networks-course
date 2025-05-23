#include <iostream>
#include <fstream>
#include <map>
#include <string>
#include <format>
#include <nlohmann/json.hpp>

using std::map;
using std::string;

using json = nlohmann::json;

struct RouteInfo {
    int cost;
    string next_hop;
};

map<string, map<string, int>> network;
map<string, map<string, RouteInfo>> routing_tables;

constexpr int INFINITY_COST = 16;

void initialize_routing_tables() {
    for (const auto& [router, neighbors] : network) {
        routing_tables[router][router] = {0, router};
        for (const auto& [neighbor, cost] : neighbors) {
            routing_tables[router][neighbor] = {cost, neighbor};
        }
    }
}

bool update_routing_tables() {
    bool updated = false;
    auto new_tables = routing_tables;

    for (const auto& [router, neighbors] : network) {
        for (const auto& [neighbor, _] : neighbors) {
            const auto& neighbor_table = routing_tables[neighbor];

            for (const auto& [dest, route] : neighbor_table) {
                if (dest == router) continue;
                int new_cost = std::min(route.cost + 1, INFINITY_COST);

                auto& entry = new_tables[router][dest];
                if (!routing_tables[router].count(dest) || new_cost < entry.cost) {
                    entry = {new_cost, neighbor};
                    updated = true;
                }
            }
        }
    }

    routing_tables = std::move(new_tables);
    return updated;
}

void print_routing_table(const map<string, map<string, RouteInfo>>& routes, int iteration = -1) {
    for (const auto& [router, table] : routes) {
        if (iteration >= 0)
            std::cout << "Iteration " << iteration << " routing table for " << router << ":\n";
        else
            std::cout << "Final routing table for " << router << ":\n";

        std::cout << std::format("{:<15} {:<15} {:<15} {:<6}\n", "Source IP", "Destination IP", "Next Hop", "Metric");
        for (const auto& [dest, info] : table) {
            std::cout << std::format("{:<15} {:<15} {:<15} {:<6}\n",
                router, dest, info.next_hop, info.cost);
        }
        std::cout << std::endl;
    }
}

int main() {
    try {
        std::ifstream config_file("../rip_config.json");
        if (!config_file.is_open()) {
            std::cerr << "Error: Failed to open configuration file.\n";
            return 1;
        }

        json config_json;
        config_file >> config_json;

        for (auto& [router, neighbors] : config_json.items()) {
            for (auto& neighbor : neighbors) {
                network[router][neighbor] = 1;
                network[neighbor][router] = 1;
            }
        }

        initialize_routing_tables();

        int iteration = 0;
        while (update_routing_tables() && iteration < 20) {
            iteration++;
            if (iteration <= 2 || iteration % 5 == 0)
                print_routing_table(routing_tables, iteration);
        }

        std::cout << "=== Final Routing Tables ===\n";
        print_routing_table(routing_tables);
    } catch (const std::exception& e) {
        std::cerr << "Exception: " << e.what() << "\n";
        return 1;
    }

    return 0;
}
