#include <stdlib.h>

int main() {
	system("ifconfig | grep \"netmask\" | awk \'{print \"IP:\", $2, \"Mask:\", $4}\'");
}
